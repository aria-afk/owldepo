create table items (
    id text primary key
);

create table item_entries (
    uuid uuid primary key default gen_random_uuid(),
    item_id text not null references items (id),
    time timestamp not null
);

create table item_entry_info (
    uuid uuid primary key default gen_random_uuid(),
    item_entry_uuid uuid not null references item_entries (uuid),
    seller_id text not null,
    quantity int not null,
    price float not null
);

-- basically k/v store of what images have already been injected
-- imagekey is defined in scrapper as: {taskId}~{screenShotTimestamp}~{filename}.png
create table images_ingested (
    imagekey text primary key
);

create table failed_images (
    imagekey text primary key
);
