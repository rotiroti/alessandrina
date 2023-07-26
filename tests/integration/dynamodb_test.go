package integration

// func runDebug(ctx context.Context) error {
// 	options := []func(*config.LoadOptions) error{}

// 	resolver := aws.EndpointResolverWithOptionsFunc(
// 		func(_, region string, options ...interface{}) (aws.Endpoint, error) {
// 			endpoint := os.Getenv("AWS_ENDPOINT_DEBUG")
// 			if endpoint != "" {
// 				return aws.Endpoint{
// 					PartitionID:   "aws",
// 					URL:           endpoint,
// 					SigningRegion: region,
// 				}, nil
// 			}

// 			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
// 		})

// 	// Enable debug logging to see the HTTP requests and responses bodies.
// 	var logMode aws.ClientLogMode

// 	if os.Getenv("AWS_CLIENT_DEBUG") != "" {
// 		logMode |= aws.LogRequestWithBody | aws.LogResponseWithBody
// 	}

// 	options = append(options,
// 		config.WithEndpointResolverWithOptions(resolver),
// 		config.WithClientLogMode(logMode),
// 	)

// 	cfg, err := config.LoadDefaultConfig(ctx, options...)
// 	if err != nil {
// 		return err
// 	}

// 	store, err := ddb.Connect(ctx, os.Getenv("TABLE_NAME"), ddb.WithAWSConfig(cfg))
// 	if err != nil {
// 		return err
// 	}

// 	bookCore := domain.NewBookCore(store)
// 	handler := web.NewAPIGatewayV2Handler(bookCore)
// 	lambda.Start(handler.CreateBook)

// 	return nil
// }

// // Config is the configuration used to create a new DynamoDB client.
// type Config struct {
// 	// ClientLog enables debug logging to see the HTTP requests and responses bodies.
// 	ClientLog string

// 	// Endpoint is the URL of the DynamoDB endpoint.
// 	Endpoint string
// }

// func parse(ctx context.Context, conf Config) ([]func(*config.LoadOptions) error, error) {
// 	options := []func(*config.LoadOptions) error{}

// 	if conf.TableName == "" {
// 		return options, fmt.Errorf("parse: %w", ErrMissingTableName)
// 	}

// 	// Define a custom endpoint resolver to use a local DynamoDB instance.
// 	// This is useful for local development, for example with the SAM CLI and LocalStack.
// 	resolver := aws.EndpointResolverWithOptionsFunc(
// 		func(_, region string, options ...interface{}) (aws.Endpoint, error) {
// 			endpoint := conf.Endpoint
// 			if endpoint != "" {
// 				return aws.Endpoint{
// 					PartitionID:   "aws",
// 					URL:           endpoint,
// 					SigningRegion: region,
// 				}, nil
// 			}

// 			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
// 		})

// 	// Enable debug logging to see the HTTP requests and responses bodies.
// 	var logMode aws.ClientLogMode

// 	if conf.ClientLog == "true" {
// 		logMode |= aws.LogRequestWithBody | aws.LogResponseWithBody
// 	}

// 	options = append(options,
// 		config.WithEndpointResolverWithOptions(resolver),
// 		config.WithClientLogMode(logMode),
// 	)

// 	return options, nil
// }

// func skipLocalstack(t *testing.T) {
// 	t.Helper()
// 	if os.Getenv("LOCALSTACK") == "" {
// 		t.Skip("skipping integration tests, set environment variable LOCALSTACK")
// 	}
// }

// func setupDB(t *testing.T) (*ddb.Store, uuid.UUID) {
// 	t.Helper()

// 	conf := ddb.Config{
// 		TableName: os.Getenv("TABLE_NAME"),
// 		Endpoint:  os.Getenv("AWS_ENDPOINT"),
// 		ClientLog: os.Getenv("AWS_CLIENT_DEBUG"),
// 	}

// 	store, err := ddb.NewStore(context.Background(), conf)
// 	require.NoError(t, err)

// 	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1cabd")

// 	return store, bookID
// }

// func TestDynamodbFlow(t *testing.T) {
// 	// Skip the DynamoDB test if the LOCALSTACK environment variable is not set
// 	skipLocalstack(t)

// 	ctx := context.Background()
// 	store, expectedBookID := setupDB(t)
// 	expectedBook := domain.Book{
// 		ID:        expectedBookID,
// 		Title:     "The Lord of the Rings",
// 		Authors:   "J.R.R. Tolkien",
// 		Pages:     1178,
// 		Publisher: "George Allen & Unwin",
// 		ISBN:      "978-0-261-10235-4",
// 	}

// 	t.Run("Save", func(t *testing.T) {
// 		err := store.Save(ctx, expectedBook)
// 		require.NoError(t, err)
// 	})

// 	t.Run("SaveTableNotExist", func(t *testing.T) {
// 		conf := ddb.Config{TableName: "not-exist"}
// 		store, err := ddb.NewStore(ctx, conf)

// 		require.NoError(t, err)

// 		err = store.Save(ctx, domain.Book{})
// 		require.Error(t, err)
// 	})

// 	t.Run("FindOne", func(t *testing.T) {
// 		err := store.Save(ctx, expectedBook)

// 		require.NoError(t, err)

// 		book, err := store.FindOne(ctx, expectedBookID)

// 		require.NoError(t, err)
// 		require.Equal(t, expectedBook, book)
// 	})

// 	t.Run("FindOneBookNotFound", func(t *testing.T) {
// 		_, err := store.FindOne(ctx, uuid.MustParse("01234567-0123-0123-0123-0123456789ab"))

// 		require.Error(t, err)
// 	})

// 	t.Run("Delete", func(t *testing.T) {
// 		err := store.Save(ctx, expectedBook)

// 		require.NoError(t, err)

// 		err = store.Delete(ctx, expectedBookID)

// 		require.NoError(t, err)
// 	})

// 	t.Run("DeleteBookNotFound", func(t *testing.T) {
// 		err := store.Delete(ctx, uuid.MustParse("01234567-0123-0123-0123-0123456789ab"))

// 		require.Error(t, err)
// 	})

// 	t.Run("FindAllBooks", func(t *testing.T) {
// 		newBooksLen := 3
// 		newBooks := make([]domain.Book, newBooksLen)

// 		for i := 0; i < newBooksLen; i++ {
// 			book := domain.Book{
// 				ID:        uuid.New(),
// 				Title:     gofakeit.BookTitle(),
// 				Authors:   gofakeit.BookAuthor(),
// 				Publisher: gofakeit.Company(),
// 				Pages:     gofakeit.Number(100, 1200),
// 			}

// 			// Save a new book in the database
// 			err := store.Save(ctx, book)
// 			require.NoError(t, err)

// 			newBooks[i] = book
// 		}

// 		books, err := store.FindAll(ctx)

// 		// Assert the expected output
// 		require.NoError(t, err)
// 		require.LessOrEqual(t, newBooksLen, len(books))
// 	})
// }
