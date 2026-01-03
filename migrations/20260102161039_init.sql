-- +goose Up
-- +goose StatementBegin

create table subscriptions (
    subscription_id bigint primary key generated always as identity,
    user_id bigint not null,
    url text not null,
    last_polled timestamptz,
    created_at timestamptz not null default now()
);

create table items (
    item_id text primary key,
    subscription_id bigint not null,
    created_at timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table items;
drop table subscriptions;
-- +goose StatementEnd
