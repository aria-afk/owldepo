-- $1 name of item example "Scroll For Helm For Dex 60%"
-- $2 amount of historical entries to fetch 
-- NOTE: Can likely index for this to be a bit faster.
SELECT JSONB_AGG(JSONB_BUILD_OBJECT(
        'time', ie.time,
        'seller_id', item_entry_info.seller_id, 
        'quantity', item_entry_info.quantity, 
        'price', item_entry_info.price
    )) 
FROM items
LEFT JOIN LATERAL (
    SELECT item_entries.uuid, item_entries.time 
    FROM item_entries
    WHERE item_entries.item_id = items.id 
    ORDER BY item_entries.time DESC 
    LIMIT $2
) ie ON TRUE
LEFT JOIN item_entry_info ON item_entry_uuid = ie.uuid
WHERE items.id = $1
