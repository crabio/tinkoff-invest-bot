# To add a new cell, type '# %%'
# To add a new markdown cell, type '# %% [markdown]'
# %% [markdown]
# ### How does it work:
# * import required libraries;
# * get data from Forbes. Link to request for top 2000 companies;
# * filtering data fields for each company;
# * create data frame;
# * save as parquet file

# %%
import pandas as pd
import requests
import os


# %%
# Request settings
params = {
    "limit": 2000
}
URL = "https://www.forbes.com/forbesapi/org/global2000/2020/position/true.json"


# %%
# Get data from Forbes website
response = requests.get(URL, params=params)
print(response.status_code)


# %%
# List of organization data
list_of_organization_data = []

# Get json from response message
original_json = response.json()

# Get only list of all organizations
organizations_json = original_json["organizationList"]["organizationsLists"]

# Loop for filtering data about ecach organization
for organization in organizations_json:
    filtered_organization = {}
    filtered_organization["Company"] = organization["organizationName"]
    filtered_organization["Country"] = organization["country"]
    filtered_organization["Industry"] = organization["industry"]
    filtered_organization["Sales"] = organization["revenue"] * 1000
    filtered_organization["Profits"] = organization["profits"] * 1000
    filtered_organization["Assets"] = organization["assets"] * 1000
    filtered_organization["Market Value"] = organization["marketValue"] * 1000

    # Add filtered data to list
    list_of_organization_data.append(filtered_organization)

# %%
# Create DataFrame from list of oraganizations dictionaries
df = pd.DataFrame(list_of_organization_data)

# Clear
df.dropna(axis=1, how='all').dropna(axis=0, how='all')

# %%
# File save settings
save_dir = "../data"
file_name = "companies_rank.csv"

# Create "data" directory if it isn't exist
if not os.path.exists(save_dir):
    os.makedirs(save_dir)

# Save to CSV
df.to_csv('{}/{}'.format(save_dir, file_name))
