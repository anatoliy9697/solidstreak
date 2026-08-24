CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
	tg_id BIGINT UNIQUE NOT NULL,
	tg_username VARCHAR(32) UNIQUE NOT NULL,
	tg_first_name VARCHAR(64),
	tg_last_name VARCHAR(64),
	tg_lang_code VARCHAR(3) NOT NULL,
	tg_is_bot BOOLEAN NOT NULL,
	lang_code VARCHAR(3) NOT NULL DEFAULT 'en',
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE tg_chats (
	id BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
	tg_id BIGINT UNIQUE NOT NULL,
	user_id BIGINT NOT NULL REFERENCES users(id),
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE habits (
	id BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	archived BOOLEAN NOT NULL DEFAULT FALSE,
	creator_id BIGINT NOT NULL REFERENCES users(id),
	title VARCHAR(256) NOT NULL,
	description TEXT,
	color VARCHAR(32) NOT NULL DEFAULT 'green' CHECK (color IN ('red', 'orange', 'yellow', 'lime', 'green', 'blue', 'purple')),
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE user_habits (
	active BOOLEAN NOT NULL DEFAULT TRUE,
	user_id BIGINT NOT NULL REFERENCES users(id),
	habit_id BIGINT NOT NULL REFERENCES habits(id),
	is_public BOOLEAN NOT NULL DEFAULT FALSE,
	PRIMARY KEY (user_id, habit_id)
);

CREATE TABLE user_habit_checks (
	user_id BIGINT NOT NULL REFERENCES users(id),
	habit_id BIGINT NOT NULL REFERENCES habits(id),
	check_date DATE NOT NULL,
	completed BOOLEAN NOT NULL DEFAULT FALSE,
	checked_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	PRIMARY KEY (user_id, habit_id, check_date)
) PARTITION BY RANGE (check_date);

CREATE TABLE user_habit_checks_y2025m09 PARTITION OF user_habit_checks
FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');

CREATE TABLE user_habit_checks_y2025m10 PARTITION OF user_habit_checks
FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE user_habit_checks_y2025m11 PARTITION OF user_habit_checks
FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');

CREATE TABLE user_habit_checks_y2025m12 PARTITION OF user_habit_checks
FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

CREATE TABLE user_habit_checks_y2026m01 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

CREATE TABLE user_habit_checks_y2026m02 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE user_habit_checks_y2026m03 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

CREATE TABLE user_habit_checks_y2026m04 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');

CREATE TABLE user_habit_checks_y2026m05 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');

CREATE TABLE user_habit_checks_y2026m06 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');

CREATE TABLE user_habit_checks_y2026m07 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');

CREATE TABLE user_habit_checks_y2026m08 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');

CREATE TABLE user_habit_checks_y2026m09 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');

CREATE TABLE user_habit_checks_y2026m10 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');

CREATE TABLE user_habit_checks_y2026m11 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');

CREATE TABLE user_habit_checks_y2026m12 PARTITION OF user_habit_checks
FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

-- Next release --

ALTER TABLE users ADD COLUMN lang_code VARCHAR(3) NOT NULL DEFAULT 'en';
UPDATE users SET lang_code = tg_lang_code;

users_habits -> user_habits

CREATE TABLE user_subscriptions (
	id BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	user_id INTEGER NOT NULL REFERENCES users(id),
	plan_code VARCHAR(32) NOT NULL CHECK (plan_code IN ('basic', 'premium')),
	start_dt TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	finish_dt TIMESTAMP WITHOUT TIME ZONE,
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

ALTER TABLE habits
ADD CONSTRAINT habits_color_check
CHECK (color IN (
    'red',
    'orange',
    'yellow',
    'lime',
    'green',
    'blue',
    'purple'
));

CREATE TABLE invoices (
	active BOOLEAN NOT NULL DEFAULT TRUE,
	uuid VARCHAR(64) NOT NULL,
	status VARCHAR(32) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'expired')),
	currency VARCHAR(32) NOT NULL CHECK (currency IN ('XTR')),
	amount BIGINT NOT NULL,
	user_id INTEGER NOT NULL REFERENCES users(id),
	tg_message_id BIGINT NOT NULL,
	tg_payment_charge_id VARCHAR(256),
	expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE subscription_events (
	id BIGSERIAL PRIMARY KEY UNIQUE NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	type VARCHAR(32) NOT NULL CHECK (type IN ('acquisition')),
	status VARCHAR(32) NOT NULL CHECK (status IN ('in_progress', 'completed', 'payment_timeout')),
	subscription_origin VARCHAR(32) CHECK (subscription_origin IN ('purchase')),
	subscription_plan_code VARCHAR(32) NOT NULL CHECK (subscription_plan_code IN ('basic', 'premium')),
	subscription_period_unit VARCHAR(32) NOT NULL CHECK (subscription_period_unit IN ('month', 'year', 'lifetime')),
	subscription_period_count BIGINT NOT NULL,
	user_id BIGINT NOT NULL REFERENCES users(id),
	subscription_id BIGINT REFERENCES user_subscriptions(id),
	invoice_uuid VARCHAR(64),
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);