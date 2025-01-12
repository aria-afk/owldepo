-- $1 name of item example "Scroll For Helm For Dex 60%"
-- $2 amount of historical entries to fetch 
-- NOTE: This can likely be optimized to be faster
SELECT JSON_BUILD_OBJECT('entries', JSON_AGG(JSON_BUILD_OBJECT(
        'time', ie.time,
        'seller_id', item_entry_info.seller_id, 
        'quantity', item_entry_info.quantity, 
        'price', item_entry_info.price
    ) ORDER BY ie.time DESC),
    'p0', MIN(item_entry_info.price),
    'p25', (SELECT PERCENTILE_DISC(0.25) WITHIN GROUP (order by item_entry_info.price)),
    'p50', (SELECT PERCENTILE_DISC(0.50) WITHIN GROUP (order by item_entry_info.price)),
    'p75', (SELECT PERCENTILE_DISC(0.75) WITHIN GROUP (order by item_entry_info.price)),
    'p100', MAX(item_entry_info.price)
) 
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
