import os
from peewee import *

script_location = os.path.dirname(os.path.realpath(__file__))
db = SqliteDatabase(os.path.join(script_location, 'database.sqlite3'))

class BaseModel(Model):
    class Meta:
        database = db

class User(BaseModel):
    username = CharField(unique=True)
    password = CharField()
    max_bytes = IntegerField()
    used_bytes = IntegerField(default=0)

class Device(BaseModel):
    mac_address = FixedCharField(unique=True, max_length=17)
    user = ForeignKeyField(User, backref='devices')

def connect():
    db.connect()
    db.create_tables([User, Device])
