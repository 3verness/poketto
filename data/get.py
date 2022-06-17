import requests
from bs4 import BeautifulSoup
from rich.progress import track
import time

base_url = "https://mikanani.me/Home/Classic/"
results = []

for i in track(range(100)):
    try:
        resp = requests.get(base_url+str(i))
        raw = BeautifulSoup(resp.text, "lxml")
        l = raw.select("#sk-body > table > tbody > tr > td:nth-child(3) > a.magnet-link-wrap")
        for t in l:
            r = t.get_text()
            results.append(r)
    except:
        break
    # time.sleep(1)

results = [i+"\n" for i in results]
with open("out.txt","w", encoding="utf-8") as f:
    f.writelines(results)