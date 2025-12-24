import random
import string
import time
import uuid
from datetime import datetime as dt

str = uuid.uuid4() 
while True:
    print(f"{dt.now()} {str}")
    time.sleep(5)
