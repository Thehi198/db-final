import pandas as pd

# Define file paths
input_json = 'News_Category_Dataset_v3.json'
output_csv = 'News_Category_Sample_100.csv'

# Load only the first 100 records
df = pd.read_json(input_json, lines=True)
df_sample = df.head(100)

# Save to CSV
df_sample.to_csv(output_csv, index=False)

print(f"Saved first 100 records to '{output_csv}'.")
