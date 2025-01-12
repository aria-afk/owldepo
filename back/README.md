# back

## Starting server

**For local development**

If you dont have air

`go install github.com/air-verse/air@latest`

`air serve.go`

## API DOC (tentative)

`/api/v1/items/recent-by-name?item={itemName}&limit={amountToRetrieve}`

**response example**

`http://localhost:8080/api/v1/items/recent-by-name?item=Scroll%20for%20Earring%20for%20INT%2060&limit=10`

```json
{
  "entries": [
    {
      "time": "2024-12-05T17:55:12",
      "seller_id": "ScrollOre",
      "quantity": 1,
      "price": 1444444
    },
    {
      "time": "2024-11-04T15:11:17",
      "seller_id": "adasik",
      "quantity": 1,
      "price": 900000
    },
    {
      "time": "2024-11-04T15:11:14",
      "seller_id": "Appo.Juice",
      "quantity": 1,
      "price": 999999
    },
    {
      "time": "2024-11-04T15:11:11",
      "seller_id": "Goodbye",
      "quantity": 1,
      "price": 888999
    },
    {
      "time": "2024-11-04T15:11:08",
      "seller_id": "amapilot",
      "quantity": 1,
      "price": 999999
    },
    {
      "time": "2024-11-04T15:11:05",
      "seller_id": "FingerHol.",
      "quantity": 1,
      "price": 999999
    },
    {
      "time": "2024-11-04T15:11:01",
      "seller_id": "Brunor",
      "quantity": 1,
      "price": 944444
    },
    {
      "time": "2024-09-12T02:16:09",
      "seller_id": "Imoeto",
      "quantity": 1,
      "price": 1222222
    },
    {
      "time": "2024-09-12T02:16:06",
      "seller_id": "M21ALI",
      "quantity": 1,
      "price": 999999
    },
    {
      "time": "2024-08-24T17:46:06",
      "seller_id": "AlyaKujou",
      "quantity": 1,
      "price": 888888
    }
  ],
  "p0": 888888,
  "p25": 900000,
  "p50": 999999,
  "p75": 999999,
  "p100": 1444444
}
```
