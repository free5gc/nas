#!/usr/bin/python3

import zipfile
import docx
import csv
from docx.oxml import CT_P
from docx.text.paragraph import Paragraph
from docx.oxml.table import CT_Tbl
from docx.table import _Cell, Table

with zipfile.ZipFile('24501-f70.zip') as zf:
    with zf.open('24501-f70.docx') as inf:
        with open('spec.csv', 'w', encoding='utf8', newline='') as outf:
            csvw = csv.writer(outf)
            doc = docx.Document(inf)
            parent = doc._body._body
            prev = None
            for c in parent.iterchildren():
                if isinstance(c, CT_P):
                    prev = Paragraph(c, parent)
                elif isinstance(c, CT_Tbl):
                    tab = Table(c, parent)
                    if isinstance(prev, Paragraph):
                        csvw.writerow([prev.text])
                        for row in tab.rows:
                            crow = []
                            for cell in row.cells:
                                crow.append(cell.text)
                            csvw.writerow(crow)
                    prev = tab
