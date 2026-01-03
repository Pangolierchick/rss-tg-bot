-- +goose Up
-- +goose StatementBegin

create table feeds (
    feed_id           bigint primary key generated always as identity,
    url               text not null unique,
    etag              text,
    last_modified     timestamptz,
    last_fetched_at   timestamptz not null default 'epoch',
    created_at        timestamptz not null default now()
);

create table subscribers (
    subscriber_id     bigint primary key generated always as identity,
    tg_chat_id        bigint not null unique,
    created_at        timestamptz not null default now()
);

create table subscriptions (
    feed_id        bigint not null references feeds(feed_id) on delete cascade,
    subscriber_id  bigint not null references subscribers(subscriber_id) on delete cascade,
    created_at     timestamptz not null default now(),
    primary key (feed_id, subscriber_id)
);

create table feed_items (
    item_id         bigint primary key generated always as identity,
    feed_id         bigint not null references feeds(feed_id) on delete cascade,
    guid            text not null,
    title           text not null,
    link            text not null,
    published_at    timestamptz,
    content_hash    bytea not null,
    created_at      timestamptz not null default now(),
    unique (feed_id, guid),
    unique (feed_id, content_hash)
);

create type delivery_status as enum ('pending', 'sent');

create table deliveries (
    delivery_id              bigint primary key generated always as identity,
    subscriber_id   bigint not null references subscribers(subscriber_id) on delete cascade,
    feed_item_id    bigint not null references feed_items(item_id) on delete cascade,
    status          delivery_status not null default 'pending',
    sent_at         timestamptz,
    created_at      timestamptz not null default now(),
    unique (subscriber_id, feed_item_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table feed_items;
drop table subscriptions;
drop table subscribers;
drop table feeds;

drop type delivery_status;
drop table deliveries;

-- +goose StatementEnd
