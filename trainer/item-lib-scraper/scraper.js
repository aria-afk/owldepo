import puppeteer from "puppeteer";

const sleep = (ms) => new Promise(res => setTimeout(res, ms));

(async () => {
    const browser = await puppeteer.launch({ headless: true });
    const page = await browser.newPage();

    // page=139 seems to be the first page with items and not just monsters.
    await page.goto("https://maplelegends.com/lib/all?page=139");
    await sleep(100);

    const tableData = await page.$$eval(".table.text-center.table-bordered > tbody > tr", tr => {
        return tr.map((trElem, i) => {
            if (i == 0) return;
            const tds = Array.from(trElem.querySelectorAll("td"));

            return {
                libHref: tds?.[1]?.querySelector("a").href,
                name: tds?.[1]?.querySelector("a").innerText,
                type: tds?.[2]?.innerText,
            }
        });
    });
    console.log(tableData);
    await browser.close();
})();
// Headers
/* TR structure:
 * #1 image container
 * <td>
     <center>
       <a href="linkToLibHref">
         <object>-</object> info of image
       </a>
     </center>
 * </td>
 * #2 Name container
 * <td>
 *   <a href="linkToLibHref">
 *     {objectName}
 *   </a>
 * </td>
 * #3 Type container 
 * <td>{objectType}</td>
*/

