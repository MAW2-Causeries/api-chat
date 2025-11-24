package databases

var migrations = map[string]map[string]string{
	"2025-11-17-11-54-24_create_messages_table": {
		"up": `
			CREATE TABLE IF NOT EXISTS messages (
				id uuid PRIMARY KEY,
				content text,
				author_id uuid,
				channel_id uuid,
				created_at timestamp,
				updated_at timestamp,
				deleted_at timestamp
			);
		`,
		"down": `
			DROP TABLE IF EXISTS messages;
		`,
	},
}

// Migrate applies the database migrations
func Migrate() error {
	if err := createMigrationsTable(); err != nil {
		return err
	}
	for name, cmds := range migrations {
		var cnt int
		if err := Session.Query(`SELECT COUNT(*) FROM migrations WHERE name = ?`, name).Scan(&cnt); err != nil {
			return err
		}
		if cnt > 0 {
			continue
		}

		if err := Session.Query(cmds["up"]).Exec(); err != nil {
			return err
		}
		if err := Session.Query(`
			INSERT INTO migrations (name, applied_at)
			VALUES (?, toTimestamp(now()))
		`, name).Exec(); err != nil {
			return err
		}
	}
	return nil
}

func createMigrationsTable() error {
	return Session.Query(`
		CREATE TABLE IF NOT EXISTS migrations (
			name text PRIMARY KEY,
			applied_at timestamp
		);
	`).Exec()
}