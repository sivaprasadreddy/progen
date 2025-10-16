package spring_boot

type DbMigrationsConfig struct {
	pg projectGenerator
}

func NewDbMigrationsConfig(pg projectGenerator) *DbMigrationsConfig {
	return &DbMigrationsConfig{pg: pg}
}

func (d DbMigrationsConfig) generate(pc ProjectConfig) error {
	return d.createDbMigrationFiles(pc)
}

func (d DbMigrationsConfig) createDbMigrationFiles(pc ProjectConfig) error {
	templateMap := map[string]string{}

	if pc.DbMigrationTool == DbMigrationToolFlyway {
		if pc.SecuritySupport || pc.JwtSecuritySupport {
			if pc.DbType == DbPostgreSQL {
				templateMap["db/migration/flyway/V1__init_postgresql.sql"] = "db/migration/V1__init.sql"
			}
			if pc.DbType == DbMySQL {
				templateMap["db/migration/flyway/V1__init_mysql.sql"] = "db/migration/V1__init.sql"
			}
			if pc.DbType == DbMariaDB {
				templateMap["db/migration/flyway/V1__init_mariadb.sql"] = "db/migration/V1__init.sql"
			}
		} else {
			templateMap["db/migration/flyway/V1__init_empty.sql"] = "db/migration/V1__init.sql"
		}
	}

	if pc.DbMigrationTool == DbMigrationToolLiquibase {
		templateMap["db/migration/liquibase/liquibase-changelog.xml"] = "db/migration/liquibase-changelog.xml"

		if pc.SecuritySupport || pc.JwtSecuritySupport {
			if pc.DbType == DbPostgreSQL {
				templateMap["db/migration/liquibase/changelog/01-init-postgresql.xml"] = "db/migration/changelog/01-init.xml"
			}
			if pc.DbType == DbMySQL {
				templateMap["db/migration/liquibase/changelog/01-init-mysql.xml"] = "db/migration/changelog/01-init.xml"
			}
			if pc.DbType == DbMariaDB {
				templateMap["db/migration/liquibase/changelog/01-init-mariadb.xml"] = "db/migration/changelog/01-init.xml"
			}
		} else {
			templateMap["db/migration/liquibase/changelog/01-init-empty.xml"] = "db/migration/changelog/01-init.xml"
		}
	}

	for tmpl, filePath := range templateMap {
		err := d.pg.executeTemplate(pc, srcMainResourcesPath+tmpl, srcMainResourcesPath+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
