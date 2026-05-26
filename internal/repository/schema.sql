pragma foreign_keys = on;

create table if not exists feeds (
    feed_id           integer primary key autoincrement,
    url               text not null unique,
    etag              text,
    last_modified     integer,
    last_fetched_at   integer not null default 0,
    created_at        integer not null default (unixepoch())
);

create table if not exists subscribers (
    subscriber_id     integer primary key autoincrement,
    tg_chat_id        integer not null unique,
    created_at        integer not null default (unixepoch())
);

create table if not exists subscriptions (
    feed_id        integer not null references feeds(feed_id) on delete cascade,
    subscriber_id  integer not null references subscribers(subscriber_id) on delete cascade,
    created_at     integer not null default (unixepoch()),
    primary key (feed_id, subscriber_id)
);

create table if not exists feed_items (
    item_id         integer primary key autoincrement,
    feed_id         integer not null references feeds(feed_id) on delete cascade,
    guid            text not null,
    title           text not null,
    link            text not null,
    published_at    integer,
    content_hash    blob not null,
    created_at      integer not null default (unixepoch()),
    unique (feed_id, guid),
    unique (feed_id, content_hash)
);

create table if not exists deliveries (
    delivery_id    integer primary key autoincrement,
    subscriber_id  integer not null references subscribers(subscriber_id) on delete cascade,
    feed_item_id   integer not null references feed_items(item_id) on delete cascade,
    status         text not null default 'pending' check (status in ('pending', 'sent')),
    sent_at        integer,
    created_at     integer not null default (unixepoch()),
    unique (subscriber_id, feed_item_id)
);

create index if not exists idx_subscriptions_subscriber on subscriptions(subscriber_id);
create index if not exists idx_feed_items_feed_created on feed_items(feed_id, created_at);
create index if not exists idx_deliveries_status_created on deliveries(status, created_at);
