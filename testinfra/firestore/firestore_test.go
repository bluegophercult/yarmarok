package firestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kaznasho/yarmarok/testinfra"
)

func TestFirestore(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)
	instance, err := RunInstance(t)
	require.NoError(t, err)

	firestoreClient := instance.Client()

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
