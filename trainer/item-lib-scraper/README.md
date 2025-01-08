# item-lib-scrapper

There are so so many items in maple legends that can be searched via owl.

This is a simple javascript script to navigate and generate a JSON list
from the `maplelegends.com/lib` page.

## To generate the items .json file

`node scraper.js`

## File structure

```json
{
    "etc"": [
        {
            "libHref": "link to maplelegnds lib page for item",
            "name": "name of item",
            "type": "Etc"
        }
    ],
    "equip": [
        {
            "libHref": "link to maplelegnds lib page for item",
            "name": "name of item",
            "type": "Equip"
        }
    ],
    "use": [
        {
            "libHref": "link to maplelegnds lib page for item",
            "name": "name of item",
            "type": Use"
        }
    ],
    "setup": [
        {
            "libHref": "link to maplelegnds lib page for item",
            "name": "name of item",
            "type": "Setup"
        }
    ]
}
```
