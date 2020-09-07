import re
from datetime import datetime
from selenium.webdriver.common.by import By
from selenium.webdriver.remote.webelement import WebElement


def getGameObj():
    return {
        "title": "",
        "details": "",
        "img_url": "",
        "url": "",
        "price": "",
        "discount_amount": "",
        "store_id": "epic",
        "category_id": "",
        "release_date": ""
    }


def transformTplToGameObj(tpl):
    info = getGameObj()

    idx = 0
    for key in info:
        info[key] = tpl[idx]
        idx += 1

    return info


def transformGameObj(gameObj):
    vals = []

    for key in gameObj:
        vals.append(gameObj[key])

    return tuple(vals)


def parseStandardStoreData(gameElement: WebElement, storeId: str, category):
    info = getGameObj()
    info["store_id"] = storeId
    info["category_id"] = category["id"]

    # Image
    try:
        imgUrl: WebElement = gameElement.find_element(
            By.CSS_SELECTOR, category["itemImageSelector"])
        info["img_url"] = imgUrl.get_attribute("src")
    except Exception as err:
        info["img_url"] = ""
        print(
            f"Image element element not found for {storeId} | {category['name']} | {category['itemImageSelector']}")

    # Title
    try:
        title: WebElement = gameElement.find_element(
            By.CSS_SELECTOR, category["itemNameSelector"])
        info["title"] = title.text
    except Exception as err:
        info["title"] = ""
        print(
            f"Title element element not found for {storeId} | {category['name']} | {category['itemNameSelector']}")

    # Url
    try:
        href: WebElement = gameElement.get_attribute("href")

        if href is not None:
            info["url"] = href
        else:
            href: WebElement = gameElement.find_element(
                By.CSS_SELECTOR, category["itemLinkSelector"]).get_attribute("href")
            info["url"] = href
    except Exception as err:
        info["url"] = ""
        print(
            f"Url element element not found for {storeId} | {category['name']} | {category['itemLinkSelector']}")

    # Details
    if category["itemDetailsSelector"] != "":
        try:
            details: WebElement = gameElement.find_element(
                By.CSS_SELECTOR, category["itemDetailsSelector"])
            info["details"] = details.text
        except Exception as err:
            info["details"] = ""
            print(
                f"Details element element not found for {storeId} | {category['name']} | {info['title']} | {category['itemDetailsSelector']}")

    # Price
    if category["itemPriceSelector"] != "":
        try:
            price = gameElement.find_element(
                By.CSS_SELECTOR, category["itemPriceSelector"])
            info["price"] = price.text
        except:
            info["price"] = ""
            print(
                f"Price element not found for {storeId} | {category['name']} | {info['title']} | {category['itemPriceSelector']}")

    # Discount amount
    if category["itemDiscountSelector"] != "":
        try:
            discountAmount = gameElement.find_element(
                By.CSS_SELECTOR, category["itemDiscountSelector"])
            info["discount_amount"] = discountAmount.text
        except:
            info["discount_amount"] = ""
            print(
                f"Discount element not found for {storeId} | {category['name']} | {info['title']} | {category['itemDiscountSelector']}")

    # Release date
    if category["releaseDateSelector"] != "":
        try:
            releaseDate = gameElement.find_element(
                By.CSS_SELECTOR, category["releaseDateSelector"])
            info["release_date"] = parseReleaseDate(releaseDate.text)
        except:
            info["release_date"] = ""
            print(
                f"Release date element not found for {storeId} | {category['name']} | {info['title']} | {category['releaseDateSelector']}")

    return info


def parseReleaseDate(releaseDate: str):
    if releaseDate == "":
        return ""

    parsed = ""

    for fmt in ("%d %b, %Y", "%b %Y", "%B %d", "%m/%d/%y"):
        if parsed != "":
            break

        try:
            parsed = datetime.strptime(releaseDate, fmt).strftime("%Y-%m-%d")

            # Handle no day case - insert 00 as a placeholder
            if fmt == "%b %Y":
                parsed = parsed.replace("-01", "-00", len(parsed) - 3)

        except:
            parsed = ""
            # print(f"Could not parse {fmt} format", err)

    return parsed or releaseDate


def parseEpicData(gameElement: WebElement, storeId: str, category):
    info = parseStandardStoreData(gameElement, storeId, category)

    itemDetails = gameElement.text.replace("\n", " ")

    # Price
    for price in re.finditer("(€|\$)\d+?(,|\.)(\d+)?|\d+?(,|\.)(\d+)?(€|\$)|Free", itemDetails):
        info["price"] = price.group()

    # Discount
    discountAmount = re.search("-\d+%", itemDetails)

    if discountAmount is not None:
        info["discount_amount"] = discountAmount.group()

    # Release date
    releaseDate = re.search("\d+\/\d+\/\d+", itemDetails)
    if releaseDate is not None:
        info["release_date"] = parseReleaseDate(releaseDate.group())

    return info
