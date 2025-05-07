import openai
import requests
from dotenv import load_dotenv
import os


load_dotenv()
openai_api_key = os.getenv("OPENAI_API_KEY")
query_api_url = "http://localhost:3000/query"


client = openai.OpenAI(api_key=openai_api_key)
def query_vector_db(query_text, k=5, metadata_filter=None):
    response = client.embeddings.create(
        input=[query_text],
        model='text-embedding-3-small'
    )
    query_vector = response.data[0].embedding

    payload = {
        "values": query_vector,
        "k": k
    }
    if metadata_filter:
        payload["metadata_filter"] = metadata_filter

    r = requests.post(query_api_url, json=payload)
    if r.status_code == 200:
        results = r.json()
        return results
    else:
        raise Exception(f"Query failed: {r.status_code} | {r.text}")

if __name__ == "__main__":
    while True:
        query = input("Enter your search query: ")
        selected_key = "headline"

        try:
            results = query_vector_db(query, k=5)
            print("\nTop Results:")
            for idx, result in enumerate(results, 1):
                metadata = result.get("vector", {}).get("metadata", {})
                value = metadata.get(selected_key, "N/A")
                distance = result.get("distance", "N/A")
                print(f"{idx}| {value} | {distance:.4f}")
        except Exception as e:
            print(e)

