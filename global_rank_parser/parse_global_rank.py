#!/usr/bin/env python
# coding: utf-8

# In[1]:


from selenium.common.exceptions import TimeoutException
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
from selenium import webdriver
from time import sleep
import pandas as pd
import os

# In[2]:


# In[12]:


CHROMEDRIVER_PATH = './../chromedriver'
URL = 'https://www.forbes.com/global2000/'

# Prepare DataFrame
df_list = []

# Init browser instance
browser = webdriver.Chrome(CHROMEDRIVER_PATH)

# Get page
browser.get(URL)
# Reload page
browser.refresh()

delay = 3  # seconds
pages_count = 5
parse_timeout_max = 10
df = pd.DataFrame()
try:
    last_id = 0
    for page_i in range(pages_count):
        # Wait table loading
        WebDriverWait(browser, delay).until(
            EC.presence_of_element_located((By.XPATH, '//table//tr//td')))
        # Get table root
        table_element = WebDriverWait(browser, delay).until(
            EC.presence_of_element_located((By.XPATH, '//table')))

        # Parse table
        table_html_string = table_element.get_attribute('outerHTML')

        # Parse table untill timeout or new data
        parse_timeout = parse_timeout_max

        while df.empty or (last_id == df.iloc[-1]['Rank']):
            # Parse table
            df = pd.read_html(table_html_string)[0]

            # Check timeout
            parse_timeout -= 1
            if parse_timeout == 0:
                raise RuntimeError("No new data for parsing.")

            # Wait reloading
            sleep(1)

        print("Last parsed rank: %d" % last_id)
        # Save new id
        last_id = df.iloc[-1]['Rank']

        # Add parsed table
        df_list.append(pd.read_html(table_html_string)[0])

        # Click NEXT button
        next_btn = browser.find_element_by_xpath(
            "//a[@ng-switch-when='next']").click()

except TimeoutException:
    print("Loading took too much time!")

# Close browser
browser.close()
browser.quit()

print("Parsed %d lists" % len(df_list))


# In[13]:


# Concat
df = pd.concat(df_list)
# Clear
df = df.dropna(axis=1, how='all').dropna(axis=0, how='all')
df


# In[49]:


# Save to parquet
df.to_parquet("../data/companies_rank.parquet")
