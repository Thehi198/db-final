import pandas as pd
import openai
import requests
import os
from dotenv import load_dotenv
import tqdm
load_dotenv()


# --- Config ---
csv_path = 'News_Category_Sample_100.csv'
column_name = 'headline'  # Choose the column to embed
openai.api_key = os.getenv('OPENAI_API_KEY')
embed_model = 'text-embedding-3-small'
insert_api_url = 'http://localhost:3000/insert'

df = pd.read_csv(csv_path)
texts = df[column_name].dropna().astype(str).tolist()

client = openai.OpenAI()

for text in tqdm(texts):

    response = client.embeddings.create(
        input=[text],
        model=embed_model
    )
    embedding = response.data[0].embedding

    payload = {
        "values": embedding,
        "metadata": {column_name: text}
    }

    r = requests.post(insert_api_url, json=payload)
    if r.status_code != 200:
        print(f"Insert failed: {text[:50]}... | Status: {r.status_code} | Response: {r.text}")
