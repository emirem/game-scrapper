import json
from time import sleep
import MySQLdb
from selenium import webdriver
import undetected_chromedriver.v2 as uc
from selenium.webdriver.common.by import By
from db import getDBConnection, getInsertQuery
from selenium.webdriver.remote.webelement import WebElement
from storeParser import transformGameObj, parseStandardStoreData, parseEpicData


def getDriver():
    options = webdriver.ChromeOptions()
    options.headless = True
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-gpu")
    options.add_argument("disable-infobars")
    options.add_argument("--disable-dev-shm-usage")

    driver: webdriver.Chrome = uc.Chrome(options=options)

    return driver


def saveToJSON(data, fileName: str):
    sanitizedName = fileName.lower().replace(" ", "_")

    with open(f"{sanitizedName}.json", 'w') as f:
        json.dump(data, f)


def scrape(driver, storeId: str, storeCategory):
    data = []
    driver.get(storeCategory["url"])

    # Arbitrary.
    sleep(1)

    try:
        elems: list[WebElement] = driver.find_elements(
            By.CSS_SELECTOR, storeCategory["listItemSelector"])

        print(f"{len(elems)} found for {storeId} - {storeCategory['id']}.")

        for elem in elems:
            info = {}

            if storeId == "epic":
                info = parseEpicData(elem, storeId, storeCategory)
            else:
                info = parseStandardStoreData(elem, storeId, storeCategory)

            if info["title"] != "":
                data.append(transformGameObj(info))
    except Exception as err:
        print("Element query failed.", err)
        return []

    return data


def scrapeStores():
    config = json.load(open("config.json"))
    # Browser
    driver = getDriver()

    # DB
    conn = getDBConnection()
    query = getInsertQuery()

    for store in config:
        for category in store["categories"]:
            categoryData = scrape(driver, store["id"], category)

            try:
                cursor: MySQLdb.cursors.Cursor = conn.cursor()
                cursor.executemany(query, categoryData)
                cursor.close()
                conn.commit()
                # saveToJSON(categoryData, f"{store['id']}_{category['id']}")
            except Exception as err:
                print("Query execute failed.", err)
                print("categoryData", categoryData)
                cursor.close()

    print("Done.")
    driver.close()
    conn.close()


if __name__ == "__main__":
    scrapeStores()