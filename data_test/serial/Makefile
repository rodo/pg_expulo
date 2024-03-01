.PHONY: test, reset


test:
	export PGPASSWORD=$(PGDSTPASSWORD)
	pg_prove -v -h $(PGDSTHOST) -p $(PGDSTPORT) -U $(PGDSTUSER) -d $(PGDSTDATABASE) source_tests.sql


reset:
	export PGPASSWORD=$(PGSRCPASSWORD)
	psql -h $(PGSRCHOST) -p $(PGSRCPORT) -U $(PGSRCUSER) -d $(PGSRCDATABASE) -f tables.sql
	psql -h $(PGSRCHOST) -p $(PGSRCPORT) -U $(PGSRCUSER) -d $(PGSRCDATABASE) -f data.sql

	export PGPASSWORD=$(PGDSTPASSWORD)
	psql -h $(PGDSTHOST) -p $(PGDSTPORT) -U $(PGDSTUSER) -d $(PGDSTDATABASE) -f tables.sql
