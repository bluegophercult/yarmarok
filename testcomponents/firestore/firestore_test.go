package firestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirestore(t *testing.T) {
	instance, err := RunInstance(t)
	require.NoError(t, err)

	firestoreClient, err := instance.Client()
	require.NoError(t, err)

	type object struct {
		Field string `firestore:"test"`
	}

	input := object{
		Field: "test",
	}

	_, err = firestoreClient.Collection("test").Doc("test").Set(context.Background(), input)

	require.NoError(t, err)

	doc, err := firestoreClient.Collection("test").Doc("test").Get(context.Background())
	require.NoError(t, err)

	var result object
	err = doc.DataTo(&result)
	require.NoError(t, err)

	require.Equal(t, input, result)

}
