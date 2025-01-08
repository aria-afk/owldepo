import puppeteer from "puppeteer";
import fs from "fs";

const sleep = (ms) => new Promise(res => setTimeout(res, ms));

(async () => {
    const browser = await puppeteer.launch({ headless: true });
    const page = await browser.newPage();

    async function parseByPage(pageNumber) {
        await page.goto(`https://maplelegends.com/lib/all?page=${pageNumber}`);
        await sleep(200);
        return await page.$$eval(".table.text-center.table-bordered > tbody > tr", tr => {
            return tr.map((trElem, i) => {
                if (i == 0) return;
                const tds = Array.from(trElem.querySelectorAll("td"));

                return {
                    // NOTE: we can eventaully use the href to get description data 
                    libHref: tds?.[1]?.querySelector("a").href,
                    name: tds?.[1]?.querySelector("a").innerText,
                    type: tds?.[2]?.innerText,
                }
            });
        });
    }
    const allPageData = [];
    // NOTE: this seems to be the range with all the non npc/cash/monster types 
    // however may be better to just fetch everything and filter after.
    let startingPage = 139;
    const endPage = 885;
    while (startingPage <= endPage) {
        process.stdout.write(`scraping pages done [ ${startingPage} ] / [ ${endPage} ]`);
        process.stdout.write(`\r`);
        allPageData.push(...(await parseByPage(startingPage)));
        startingPage++;
    }

    const resJson = {
        "etc": [],
        "equip": [],
        "use": [],
        "setup": [],
    };

    for (let i = 0; i < allPageData.length; i++) {
        const pd = allPageData[i];
        switch (pd?.type) {
            case "Etc":
                resJson.etc.push(pd);
            case "Equip":
                resJson.equip.push(pd);
            case "Use":
                resJson.use.push(pd);
            case "Setup":
                resJson.setup.push(pd);
            default:
                continue;
        }
    }

    if (fs.existsSync("./mapleitems.json")) {
        fs.unlinkSync("./mapleitems.json");
    }
    fs.writeFileSync("./mapleitems.json", JSON.stringify(resJson));

    await browser.close();
})();
