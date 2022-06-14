package database

// type Cfg struct {
// 	URI     string
// 	Name    string
// 	Timeout time.Duration
// }

// func Open(cfg Cfg) (*mongo.Client, context.Context, context.CancelFunc, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
// 	if err != nil {
// 		cancel()
// 		return nil, nil, nil, fmt.Errorf("opening db: %w", err)
// 	}

// 	return client, ctx, cancel, nil
// }
