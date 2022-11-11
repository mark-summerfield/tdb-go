# TDB Overview

Text DataBase (TDB) is a plain text human readable typed database storage
format.

- [Datatypes](#datatypes)
    - [Formatting](#formatting)
- [Examples](#examples)
    - [CSV](#csv)
    - [Database](#database)
- [Libraries](#libraries)
- [BNF](#bnf)
- [Supplementary](#supplementary)
    - [Vim Support](#vim-support)
    - [TDB Logo](#tdb-logo)

## Datatypes

TDB supports the following six built-in datatypes.

|**Type**<a name="table-of-built-in-types"></a>|**Example(s)**|**Notes**|
|-----------|----------------------|--|
|`bool`     |`F` `T`||
|`bytes`    |`(20AC 65 66 48)`|There must be an even number of case-insensitive hex digits; whitespace (spaces, newlines, etc.) optional.|
|`date`     |`2022-04-01`|Basic ISO8601 YYYY-MM-DD format.|
|`datetime` |`2022-04-01T16:11:51`|ISO8601 YYYY-MM-DDTHH[:MM[:SS]] format; 1-sec resolution no timezone support.|
|`int`      |`-192` `+234` `7891409`|Standard integers with optional sign.|
|`real`     |`0.15` `0.7e-9` `2245.389`|Standard and scientific notation.|

All these types are _not null_, that is, null values (represented by `?`)
are invalid in fields whose type is one of the above.

To allow a field to be null, use the type followed by `?`, e.g., `int?` is a
nullable `int`.

Strings may not include `&`, `<` or `>`, so if they are needed, they must be
replaced by the XML/HTML escapes `&amp;`, `&lt;`, and `&gt;` respectively.
Strings respect any whitespace they contain, including newlines.

Where whitespace is allowed (or required) it may consist of one or more
spaces, tabs, or newlines in any combination.

### Formatting

A TDB file's header must always occupy its own line (i.e., end with a
newline). The rest of the file could in theory be a single line no matter
how long. In practice and for human readability it is normal to limit the
width of lines, for example, to 76, 80, or the TDB default of 96 characters.

A TDB processor is expected to provide formatting options for pretty
printing TDB files with user defined indentation, wrap width, and real
number formatting.

TDB `bytes` and ``str``s can be of any length, but nonetheless they can be
width-limited without changing their semantics.

#### Bytes

Any `bytes` value may be written with any amount of whitespace including
newlines—with all the whitespace ignored. For example:

    (AB DE 01 57) ≣ (ABDE0157)

This makes it is easy to convert a `bytes` that is too long into chunks,
e.g.,

    (20 AC 40 41 ... lots more ... FF FE)

to, say:

    (20 AC 40 41
    ... some more ...
    ... some more ...
    FF FE)

#### Strings

Because TDB strings respect any whitespace they contain they cannot be split
into chunks like `bytes`. However, TDB supports a string concatenation
operator such that:

    <This is one string> ≣ <This > & <is one > & <string>

Which means, of course, that given a long string that might not contain
newlines or whose lines are too long, we can easily split it into chunks,
e.g.,

    <Imagine this is a really long string...>

to, say:

    <Imagine > &
    <this is a > &
    <really long > &
    <string...>

Comments work the same way, but note that the comment marker must only
precede the _first_ fragment.

    #<This is a comment in one or more strings.> ≣ #<This is a > & <comment in > & <one or more> & < strings.>

## Examples

### CSV

Although widely used, the CSV format is not standardized and has a number of
problems. TDB is a standardized alternative that can distinguish fieldnames
from data rows, can handle multiline text (including text with commas and
quotes) without formality, and can store one—or more—tables in a single TDB
file.

Here's a simple CSV file:

    Date,Price,Quantity,ID,Description
    "2022-09-21",3.99,2,"CH1-A2","Chisels (pair), 1in & 1¼in"
    "2022-10-02",4.49,1,"HV2-K9","Hammer, 2lb"
    "2022-10-02",5.89,1,"SX4-D1","Eversure Sealant, 13-floz"

Here's a TDB equivalent:

    TDB1
    [PriceList Date date Price real Quantity int ID str Description str?
    %
      2022-09-21 3.99 2 <CH1-A2> <Chisels (pair), 1in &amp; 1¼in> 
      2022-10-02 4.49 1 <HV2-K9> <Hammer, 2lb> 
      2022-10-02 5.89 1 <SX4-D1> <Eversure Sealant, 13-floz> 
      2022-11-13 8.49 1 <PV7-X2> ?
    ]

Every table starts with a table name followed by one or more fields. Each
field consists of a field name and a type.

In this example the `Date`, `Price`, `Quantity`, and `ID` fields are _not
null_, but the `Description` field is a nullable `str`.

### Database

Database files aren't normally human readable and usually require
specialized tools to read and modify their contents. Yet many databases are
relatively small (both in size and number of tables), and would be more
convenient to work with if human readable. For these, TDB format provides a
viable alternative.

A TDB equivalent to a database of tables can easily be created using a
`list` of ``table``s:

    TDB1 MyApp Data
    [Customers CID int Company str Address str? Contact str? Email str
    %
        50 <Best People> <123 Somewhere> <John Doe> <j@doe.com> 
        19 <Supersuppliers> ? <Jane Doe> <jane@super.com> 
    ]
    [Invoices INUM int CID int Raised_Date date Due_Date date Paid bool
    Description str?
    %
        152 50 2022-01-17 2022-02-17 no <COD> 
        153 19 2022-01-19 2022-02-19 yes <> 
    ]
    [Items IID int INUM int Delivery_Date date Unit_Price real Quantity int
    Description str?
    %
        1839 152 2022-01-16 29.99 2 <Bales of hay> 
        1840 152 2022-01-16 5.98 3 <Straps> 
        1620 153 2022-01-19 11.5 1 <Washers (1-in)> 
    ]


## Libraries

|**Library**|**Language**|**Notes**                    |
|-----------|------------|-----------------------------|

## BNF

A TDB file consists of a mandatory header followed by an optional file-level
comment, optional imports, optional _ttype_ definitions, and then a single
mandatory `list`, `map`, or `table` (which may be empty).

    TDB          ::= 'TDB1' CUSTOM? '\n' TABLE+
    CUSTOM       ::= RWS [^\n]+ # user-defined data e.g. filetype and version
    TABLE        ::= '[' OWS TDEF OWS '%' OWS ROW* OWS ']'
    TDEF         ::= IDENFIFIER (RWS FDEF)+ # IDENFIFIER is the table name
    FDEF         ::= IDENFIFIER RWS TYPE ('?')? # IDENFIFIER is the field name
    TYPE         ::= 'bool' | 'bytes' | 'date' | 'datetime' | 'int' | 'real' | 'str'
    ROW          ::= OWS VALUE (RWS VALUE)*
    VALUE        ::= BOOL | BYTES | DATE | DATETIME | INT | NULL | REAL | STR # NULL is only valid in nullable fields
    NULL         ::= '?'
    BOOL         ::= 'F' | 'T'
    INT          ::= /[-+]?\d+/
    REAL         ::= # standard or scientific notation
    DATE         ::= /\d\d\d\d-\d\d-\d\d/ # basic ISO8601 YYYY-MM-DD format
    DATETIME     ::= /\d\d\d\d-\d\d-\d\dT\d\d(\d\d(\d\d)?)?/ # see note below
    STR          ::= STR_FRAGMENT (OWS '&' OWS STR_FRAGMENT)*
    STR_FRAGMENT ::= /[<][^<>]*?[>]/ # newlines allowed, and &amp; &lt; &gt; supported i.e., XML
    BYTES        ::= '(' (OWS [A-Fa-f0-9]{2})* OWS ')'
    IDENFIFIER   ::= /[_\p{L}]\w{0,31}/ # Must start with a letter or underscore; may not be a built-in TYPE
    OWS          ::= /[\s\n]*/
    RWS          ::= /[\s\n]+/ # in some cases RWS is actually optional

Note that a TDB file _must_ contain at least one table even if it is empty.

Note that for any given table each field name must be unique.

Note that a TDB reader (writer) _must_ be able to read (write) a plain text
`.tdb` file containing UTF-8 encoded text, and _ought_ to be able to read
and write gzipped plain text `.tdb.gz` files.

Note also that TDB readers and writers should _not_ care about the actual
file extension (apart from the `.gz` needed for gzipped files), since users
are free to use their own. For example, `data.myapp` and `data.myapp.gz`.

## Supplementary

### Vim Support

If you use the vim editor, simple color syntax highlighting is available.
Copy `tdb.vim` into your `$VIM/syntax/` folder and add these lines (or
similar) to your `.vimrc` or `.gvimrc` file:

    au BufRead,BufNewFile,BufEnter * if getline(1) =~ '^tdb ' | setlocal ft=tdb | endif
    au BufRead,BufNewFile,BufEnter *.tdb set ft=tdb|set expandtab|set tabstop=2|set softtabstop=2|set shiftwidth=2

### TDB Logo

![tdb logo](tdb.svg)

---
