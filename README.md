# Tdb Overview

Text database (Tdb) is a plain text human readable typed database storage
format.

- [Datatypes](#datatypes)
    - [Formatting](#formatting)
- [Examples](#examples)
    - [Minimal Tdb Files](#minimal-tdb-files)
    - [CSV](#csv)
    - [Database](#database)
- [Libraries](#libraries)
	- [Go](#go)
- [BNF](#bnf)
- [Supplementary](#supplementary)
    - [Vim Support](#vim-support)
    - [Tdb Logo](#tdb-logo)

## Datatypes

Tdb supports the following seven built-in datatypes.

|**Type**<a name="table-of-built-in-types"></a>|**Zero Value**|**Example(s)**|**Notes**|
|-----------|----------------------|--|
|`bool`     |`F`|`F` `T`||
|`bytes`    |`()`|`(20AC 65 66 48)`|There must be an even number of case-insensitive hex digits; whitespace (spaces, newlines, etc.) optional.|
|`date`     |`1900-01-01`   |`2022-04-01`|Basic ISO8601 YYYY-MM-DD format.|
|`datetime` |`1900-01-01T00`|`2022-04-01T16:11:51`|ISO8601 YYYY-MM-DDTHH[:MM[:SS]] format; 1-sec resolution no timezone support.|
|`int`      |`0`|`-192` `+234` `7891409`|Standard integers with optional sign.|
|`real`     |`0.0`|`0.15` `0.7e-9` `2245.389`|Standard and scientific notation.|
|`str`      |`<>`|`<Some text which may include newlines>`|For &, <, >, use \&amp;, \&lt;, \&gt; respectively.|

Nulls are not supported. However, each type has a "zero” value which may be
useful. (In some cases a null could be simulated by a sentinal value.)

Strings may not include `&`, `<` or `>`, so if they are needed, they must be
replaced by the XML/HTML escapes `&amp;`, `&lt;`, and `&gt;` respectively.
Strings respect any whitespace they contain, including newlines.

Where whitespace is allowed (or required) it may consist of one or more
spaces, tabs, or newlines in any combination.

### Formatting

A Tdb file's header must always occupy its own line (i.e., end with a
newline). The rest of the file could in theory be a single line no matter
how long. In practice and for human readability it is normal to limit the
width of lines, for example, to 76, 78, or the Tdb default of 80 characters.

A Tdb processor is expected to provide formatting options for pretty
printing Tdb files with user defined indentation, wrap width, and real
number formatting.

Tdb `bytes` and ``str``s can be of any length, but nonetheless they can be
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

Because Tdb strings respect any whitespace they contain they cannot be split
into chunks like `bytes`. However, Tdb supports a string concatenation
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

### Minimal Tdb Files

	[T f int]

This file has a single table called `T` which has a single field called `f`
of type `int`, and no rows.

	[T f int % 0]

This is like the previous table but now with one row containing the value
`0`.

	Tdb1
	[T f int % 0 #]

Again like the previous table, but now with the optional header and two
rows, the first containing the value `0`, and the second also containing the
value `0` (`#` is used to indicate the field type's “zero” value, in this
case, `0`).

### CSV

Although widely used, the CSV format is not standardized and has a number of
problems. Tdb is a standardized alternative that can distinguish fieldnames
from data rows, can handle multiline text (including text with commas and
quotes) without formality, and can store one—or more—tables in a single Tdb
file.

Here's a simple CSV file:

    Date,Price,Quantity,ID,Description
    "2022-09-21",3.99,2,"CH1-A2","Chisels (pair), 1in & 1¼in"
    "2022-10-02",4.49,1,"HV2-K9","Hammer, 2lb"
    "2022-10-02",5.89,1,"SX4-D1","Eversure Sealant, 13-floz"

Here's a Tdb equivalent:

    Tdb1
    [PriceList Date date Price real Quantity int ID str Description str
    %
      2022-09-21 3.99 2 <CH1-A2> <Chisels (pair), 1in &amp; 1¼in> 
      2022-10-02 4.49 1 <HV2-K9> <Hammer, 2lb> 
      2022-10-02 5.89 1 <SX4-D1> <Eversure Sealant, 13-floz> 
      2022-11-13 8.49 # <PV7-X2> #
    ]

Every table starts with a table name followed by one or more fields. Each
field consists of a field name and a type.

In this example the last row uses `#` to indicate the field type's “zero”
value (in this case a quantity of `0` and an empty description string).

### Database

Database files aren't normally human readable and usually require
specialized tools to read and modify their contents. Yet many databases are
relatively small (both in size and number of tables), and would be more
convenient to work with if human readable. For these, Tdb format provides a
viable alternative. For example:

    Tdb1
    [Customers CID int Company str Address str Contact str Email str
    %
        50 <Best People> <123 Somewhere> <John Doe> <j@doe.com> 
        19 <Supersuppliers> # <Jane Doe> <jane@super.com> 
    ]
    [Invoices INUM int CID int Raised_Date date Due_Date date Paid bool Description str
    %
        152 50 2022-01-17 2022-02-17 no <COD> 
        153 19 2022-01-19 2022-02-19 yes # 
    ]
    [Items IID int INUM int Delivery_Date date Unit_Price real Quantity int Description str
    %
        1839 152 2022-01-16 29.99 2 <Bales of hay> 
        1840 152 2022-01-16 5.98 3 <Straps> 
        1620 153 2022-01-19 11.5 1 <Washers (1-in)> 
    ]

In the Customers table the second customer's Address and in the Invoices
table, the second invoice's Description the `str` type's zero value (an
empty string) is used.

## Libraries

|**Library**|**Language**|**Notes**                    |
|-----------|------------|-----------------------------|

### Go

|**Tdb Type**|**Go Type**|
|------------|----------------------|
|`bool`      |`bool`|
|`bytes`     |`[]byte`|
|`date`      |`time.Time`|
|`datetime`  |`time.Time`|
|`int`       |`int`|
|`real`      |`float64`|
|`str`       |`string`|

## BNF

A Tdb file consists of an optional header followed by one or more tables.

    TDB          ::= ('Tdb1' CUSTOM? '\n')? TABLE+
    CUSTOM       ::= RWS [^\n]+ # user-defined data e.g. filetype and version
    TABLE        ::= OWS '[' OWS TDEF OWS '%' OWS ROW* OWS ']'
    TDEF         ::= IDENFIFIER (RWS FDEF)+ # IDENFIFIER is the table name
    FDEF         ::= IDENFIFIER RWS TYPE # IDENFIFIER is the field name
    TYPE         ::= 'bool' | 'bytes' | 'date' | 'datetime' | 'int' | 'real' | 'str'
    ROW          ::= OWS VALUE (RWS VALUE)*
    VALUE        ::= BOOL | BYTES | DATE | DATETIME | INT | REAL | STR 
    ZERO         ::= '#'
    BOOL         ::= 'F' | 'T' | ZERO
    INT          ::= /[-+]?\d+/ | ZERO
    REAL         ::= # standard or scientific notation or ZERO
    DATE         ::= /\d\d\d\d-\d\d-\d\d/ | ZERO # basic ISO8601 YYYY-MM-DD format
    DATETIME     ::= /\d\d\d\d-\d\d-\d\dT\d\d(\d\d(\d\d)?)?/ | ZERO
    STR          ::= STR_FRAGMENT (OWS '&' OWS STR_FRAGMENT)* | ZERO
    STR_FRAGMENT ::= /[<][^<>]*?[>]/ # newlines allowed, and &amp; &lt; &gt; supported i.e., XML
    BYTES        ::= '(' (OWS [A-Fa-f0-9]{2})* OWS ')' | ZERO
    IDENFIFIER   ::= /[_\p{L}]\w{0,31}/ # Must start with a letter or underscore; may not be a built-in TYPE
    OWS          ::= /[\s\n]*/
    RWS          ::= /[\s\n]+/ # in some cases RWS is actually optional

Note that a Tdb file _must_ contain at least one table even if it is empty,
i.e., has no rows.

Note that for any given table each field name must be unique.

Note that a Tdb reader (writer) _must_ be able to read (write) a plain text
`.tdb` file containing UTF-8 encoded text, and _ought_ to be able to read
and write gzipped plain text `.tdb.gz` files.

Note also that Tdb readers and writers should _not_ care about the actual
file extension (apart from the `.gz` needed for gzipped files), since users
are free to use their own. For example, `data.myapp` and `data.myapp.gz`.

## Supplementary

### Vim Support

If you use the vim editor, simple color syntax highlighting is available.
Copy `tdb.vim` into your `$VIM/syntax/` folder and add these lines (or
similar) to your `.vimrc` or `.gvimrc` file:

    au BufRead,BufNewFile,BufEnter * if getline(1) =~ '^Tdb1' | setlocal ft=tdb | endif
    au BufRead,BufNewFile,BufEnter *.tdb set ft=tdb|set expandtab|set tabstop=2|set softtabstop=2|set shiftwidth=2

### Tdb Logo

![tdb logo](tdb.svg)

---
