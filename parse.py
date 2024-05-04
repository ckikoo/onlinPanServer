import PyPDF2

def get_bookmarks_with_page_numbers(pdf_file_path):
    bookmarks_with_page_numbers = []

    with open(pdf_file_path, 'rb') as pdf_file:
        reader = PyPDF2.PdfReader(pdf_file)
        num_pages = len(reader.pages)
        outlines = reader.outline

        def recursive(bookmarks, parent_title='', parent_page=None):
            for bookmark in bookmarks:
                title = bookmark.title
                dest = bookmark.get('/Dest', None)
                if dest:
                    page_number = reader.get_reference(dest)[0] + 1
                    bookmarks_with_page_numbers.append((parent_title + '/' + title, page_number))
                if isinstance(bookmark, list):
                    recursive(bookmark, parent_title + '/' + title, parent_page)
        
        recursive(outlines)
    
    return num_pages, bookmarks_with_page_numbers

# 替换 'your_pdf_file.pdf' 为你的 PDF 文件路径
pdf_file_path = '1python.pdf'
total_pages, bookmarks_with_page_numbers = get_bookmarks_with_page_numbers(pdf_file_path)

print(f"Total Pages: {total_pages}\n")

for bookmark, page_number in bookmarks_with_page_numbers:
    print(f"Bookmark: {bookmark}, Page Number: {page_number}")
