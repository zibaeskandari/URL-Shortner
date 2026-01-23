## Database Tables
### USER
We would have three type of user `user`, `admin`, and `superAdmin`. We can use file based (a json for example for permissions, or we can store permissions in a table per role). However since the project is small we can go with file based permissions.
```sql
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(32) NOT NULL,
                       password_hash VARCHAR(128) NOT NULL,
                       role VARCHAR(32) NOT NULL DEFAULT 'user',
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ux_users_username_active ON users (username) WHERE deleted_at IS NULL;
```

### Token
We use paseto token to secure our APIs. we also have refresh token. We may consider a simpler approach to rotate tokens. We update the `paseto_refresh_token` on each refresh token call. And if a refresh token call uses an invalid token we revoke the user refresh token altogether. And user has to do a login. This approach will cause a normal user to be logged out from a device if they login from another device. Although this inconvenient, in cases of a compromised token, we are not providing service.
```sql

CREATE TABLE paseto_refresh_token (
    id bigserial primary key,
    user_id INTEGER NOT NULL,
    token_fingerprint VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_paseto_refresh_token_user_id ON paseto_refresh_token (user_id, token_fingerprint)
````

### URLs
```sql
CREATE TABLE urls (
      id VARCHAR(32) PRIMARY KEY,
      destination TEXT NOT NULL,
      user_id INTEGER NOT NULL,
      expires_at TIMESTAMPTZ,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX ix_urls_user_active ON urls (user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX ix_urls_expires ON urls (expires_at) WHERE expires_at IS NOT NULL AND deleted_at IS NULL;
ALTER TABLE urls ADD CONSTRAINT fk_urls_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
```

### Usage
Since we do not want to block user on resolving urls, we just insert each short url's usage in a `write ahead log` kind of table. We have another worker that will collect and aggregate this log. Indexes will slow the insertion, so we can have another worker to insert into this table and app only puts the log in channel of this worker. Also we can ignore this table altogether and feed the log to aggregation worker via channel.
```sql
CREATE TABLE usage_wal (
       id BIGSERIAL PRIMARY KEY,
       url_id VARCHAR(32) NOT NULL,
       user_agent TEXT,
       ip_address INET,
       referrer TEXT,
       aggregated BOOLEAN DEFAULT FALSE,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_usage_wal_url_time ON usage_wal (url_id, created_at DESC);
CREATE INDEX ix_usage_wal_id_aggregated ON usage_wal (id, aggregated) where aggregated = TRUE;
```

This table will keep a usage per day for each `url_id`. Since we have considered `day DATE NOT NULL` it only keeps the usage per day in `UTC` days, if we want to be able to query for all timezones we can aggregate on 30 minutes (minimum difference of each timezone) window and have 48 rows per day for each `url_id`. 

```sql
CREATE TABLE usage_log (
   url_id VARCHAR NOT NULL,
   day DATE NOT NULL,
   count_per_day INTEGER NOT NULL CHECK (count_per_day >= 0),
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   PRIMARY KEY (url_id, day)
);

ALTER TABLE usage_log ADD CONSTRAINT fk_usage_log_url FOREIGN KEY (url_id) REFERENCES urls(id) ON DELETE CASCADE;
```