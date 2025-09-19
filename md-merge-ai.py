import os
import argparse

def merge_markdown_files(input_dir, output_file, ai_context="", ai_prompt="", separator="---"):
    """
    Merges all Markdown files in a directory into a single, well-structured file,
    with an AI context and prompt at the beginning.

    Args:
        input_dir (str): The directory containing the Markdown files.
        output_file (str): The name of the output Markdown file.
        ai_context (str): A string providing context to the AI about the documents.
        ai_prompt (str): A string with specific instructions for the AI.
        separator (str): A string to use as a separator between articles.
    """
    if not os.path.exists(input_dir):
        print(f"Error: The directory '{input_dir}' does not exist.")
        return

    # Create the output file with a clear header and AI instructions
    with open(output_file, 'w', encoding='utf-8') as outfile:
        # Write the AI context section
        if ai_context:
            print(f"### AI Context\n\n{ai_context}\n\n{separator}\n\n", file=outfile)
            
        # Write the AI prompt section
        if ai_prompt:
            print(f"### AI Prompt\n\n{ai_prompt}\n\n{separator}\n\n", file=outfile)
            
        print(f"# Merged Knowledge Base for AI\n\n", file=outfile)
        
        # Get a list of all .md files in the input directory
        markdown_files = sorted([f for f in os.listdir(input_dir) if f.endswith('.md')])

        if not markdown_files:
            print(f"No Markdown files found in the directory: {input_dir}")
            return

        # Iterate through each file and append its content to the output file
        for i, filename in enumerate(markdown_files):
            filepath = os.path.join(input_dir, filename)
            
            # Use the filename as the title for the new section
            title = os.path.splitext(filename)[0]
            
            print(f"## Article: {title}\n\n", file=outfile)
            
            with open(filepath, 'r', encoding='utf-8') as infile:
                content = infile.read()
                outfile.write(content)
                
            if i < len(markdown_files) - 1:
                print(f"\n\n{separator}\n\n", file=outfile)

    print(f"Successfully merged {len(markdown_files)} files into {output_file}.")
    print("The final file is now structured for efficient AI processing.")

# --- CLI Usage ---
if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="""
        Merge multiple Markdown files from a directory into a single,
        AI-ready document. This script allows you to add a specific context
        and prompt for better AI performance.
        """,
        epilog="""
        Examples:
        
        # Merge all files in 'my_docs' and create a file named 'knowledge.md'
        python merge_markdown.py my_docs knowledge.md

        # Merge files and include a context and prompt for the AI
        python merge_markdown.py my_docs knowledge.md \\
        --context "This is a company knowledge base containing technical guides." \\
        --prompt "You are a technical support agent. Answer questions using only this document."
        """
    )
    
    parser.add_argument("input_dir", help="The directory containing the Markdown files to merge.")
    parser.add_argument("output_file", help="The name of the output Markdown file.")
    parser.add_argument("--context", dest="ai_context", default="", help="A string providing context for the AI.")
    parser.add_argument("--prompt", dest="ai_prompt", default="", help="A string with specific instructions for the AI.")
    parser.add_argument("--separator", dest="separator", default="---", help="A separator to use between articles.")

    args = parser.parse_args()

    merge_markdown_files(args.input_dir, args.output_file, args.ai_context, args.ai_prompt, args.separator)

