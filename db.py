import MySQLdb
import os
from dotenv import load_dotenv

load_dotenv()


def getDBConnection() -> MySQLdb.Connection:
    return MySQLdb.connect(
        host=os.getenv("HOST"),
        user=os.getenv("USERNAME"),
        passwd=os.getenv("PASSWORD"),
        db=os.getenv("DATABASE"),
        ssl_mode="VERIFY_IDENTITY",
        ssl={
            "ca": "/etc/ssl/certs/ca-certificates.crt"
        }
    )


def getInsertQuery():
    query = ("INSERT INTO `games` (title, details, img_url, url, price, discount_amount, store_id, category_id, release_date) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)")

    return query
