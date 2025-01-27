.PHOMY: apply

apply:
	atlas migrate apply \
	--dir "file://ent/migrate/migrations" \
	--url sqlite3://test.db

