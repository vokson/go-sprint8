package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func newTestStore(t *testing.T) ParcelStore {
	t.Helper()

	db, err := sql.Open("sqlite", DBFileName)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	return NewParcelStore(db)
}

func TestAddGetDelete(t *testing.T) {
	store := newTestStore(t)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	parcel.Number = id

	dbParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, parcel, dbParcel)

	require.NoError(t, store.Delete(id))

	_, err = store.Get(id)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSetAddress(t *testing.T) {
	store := newTestStore(t)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	dbParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, dbParcel.Address)
}

func TestSetStatus(t *testing.T) {
	store := newTestStore(t)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	newStatus := "new test status"
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	dbParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newStatus, dbParcel.Status)
}

func TestGetByClient(t *testing.T) {
	store := newTestStore(t)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)

	assert.Len(t, storedParcels, len(parcels))

	for _, parcel := range storedParcels {
		p, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		assert.Equal(t, parcel, p)
	}
}
