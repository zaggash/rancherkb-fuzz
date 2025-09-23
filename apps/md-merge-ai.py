import os
import argparse

def merge_markdown_files(input_dir, output_file, ai_context="", ai_prompt=""):
    """
    Merges all Markdown files from a directory and its subdirectories into a
    single, well-structured file, optimized for AI analysis.

    Args:
        input_dir (str): The base directory to scan for Markdown files.
        output_file (str): The name of the output Markdown file.
        ai_context (str): Contextual information for the AI.
        ai_prompt (str): Specific instructions for the AI's behavior.
    """
    if not os.path.isdir(input_dir):
        print(f"Error: The input directory '{input_dir}' does not exist or is not a directory.")
        return

    # Use a standard separator for the articles
    article_separator = "---"

    # Find all Markdown files by walking the directory tree
    markdown_files = []
    for root, dirs, files in os.walk(input_dir):
        for file in sorted(files):
            if file.endswith('.md'):
                markdown_files.append(os.path.join(root, file))

    if not markdown_files:
        print(f"No Markdown files found in '{input_dir}' or its subdirectories.")
        return

    # Write the merged content to the output file
    with open(output_file, 'w', encoding='utf-8') as outfile:
        # Write the AI context and prompt section without a special separator
        if ai_context:
            outfile.write(f"### AI Context\n\n{ai_context}\n\n")

        if ai_prompt:
            outfile.write(f"### AI Prompt\n\n{ai_prompt}\n\n")

        # Write the main document title
        outfile.write(f"# Merged Knowledge Base for AI\n\n")

        # Iterate through each file and append its content
        for i, filepath in enumerate(markdown_files):
            # Read the content of the current file
            with open(filepath, 'r', encoding='utf-8') as infile:
                content = infile.read()
            
            # Add the appropriate separator before the article, but not for the first file
            if i > 0:
                outfile.write(f"\n\n{article_separator}\n\n")

            # Use a clear heading for each file
            relative_path = os.path.relpath(filepath, input_dir)
            outfile.write(f"## Article: {relative_path}\n\n")

            # Write the file's content
            outfile.write(content)

    print(f"Successfully merged {len(markdown_files)} files into '{output_file}'.")
    print("The final file is now structured for efficient AI processing.")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="""
        Merge multiple Markdown files from a directory and its subdirectories
        into a single, AI-ready document. This script allows you to add specific
        context and instructions to optimize the output for an AI.
        """,
        epilog="""
        ---
        Examples:
        
        # Merge all files in 'docs' and its subfolders into 'knowledge.md'
        python merge_markdown.py docs knowledge.md

        # Merge files and include a context and prompt for the AI
        python merge_markdown.py docs knowledge.md \\
        --context "This is a technical knowledge base for our new software. It contains API documentation and bug reports." \\
        --prompt "You are a technical expert. Your purpose is to answer user questions using only the information provided in this document."
        """
    )
    
    parser.add_argument("input_dir", help="The directory to scan for Markdown files.")
    parser.add_argument("output_file", help="The name of the output Markdown file.")
    parser.add_argument("--context", dest="ai_context", default="", help="A string providing context for the AI.")
    parser.add_argument("--prompt", dest="ai_prompt", default="", help="A string with specific instructions for the AI.")

    args = parser.parse_args()

    merge_markdown_files(args.input_dir, args.output_file, args.ai_context, args.ai_prompt)

