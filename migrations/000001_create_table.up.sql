CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id uuid NOT NULL,
    service VARCHAR(30) NOT NULL,
    price INT NOT NULL CHECK (price >= 0),
    subscribed_on DATE NOT NULL,
    unsubscribed_on DATE CHECK (ubsubscribed_on >= subscribed_on)
);

CREATE INDEX IF NOT EXISTS idx_subscription_user_ids ON subscriptions (user_id);
