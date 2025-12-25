import random
import string
import time
import uuid
from datetime import datetime as dt
import os

str = uuid.uuid4() 
while True:
    out_path = os.path.expanduser(os.path.expandvars(os.environ.get("OUTPUT_PATH", "output.txt")))
    with open(out_path, "a", encoding="utf-8") as f:
      f.write(f"{dt.now()} {str}\n")
    time.sleep(5)