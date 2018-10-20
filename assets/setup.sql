-- Config
CREATE TABLE IF NOT EXISTS config (
    confkey VARCHAR(255),
    confvalue VARCHAR(255)
);

-- Fill Config
INSERT INTO config
    (confkey, confvalue)
SELECT 'title', 'Gorum'
WHERE
    NOT EXISTS (
        SELECT confkey FROM config WHERE confkey = 'title'
    );

-- Users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    username VARCHAR(255),
    passwordhash VARCHAR(255),
    mail VARCHAR(255),
    registered VARCHAR(255),
    avatar TEXT DEFAULT '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' NOT NULL
);