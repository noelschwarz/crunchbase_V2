# Research NoSQL Databases
This document presents the findings from the research into NoSQL databases that could be used in this project.

## Table of contents
<!-- vim-markdown-toc GFM -->

* [Data model](#data-model)
	- [Basic requirements](#basic-requirements)
* [Why NoSQL makes sense for this project?](#why-nosql-makes-sense-for-this-project)
	- [Mutable data schema](#mutable-data-schema)
* [Document-oriented databases](#document-oriented-databases)
	- [Main characteristics of document-oriented databases](#main-characteristics-of-document-oriented-databases)
* [Time-series databases](#time-series-databases)
	- [Arguments against TSDBs](#arguments-against-tsdbs)
* [MongoDB](#mongodb)
	- [Collections](#collections)
	- [Time-series collections](#time-series-collections)
	- [Indexing](#indexing)
	- [Schema versioning](#schema-versioning)

<!-- vim-markdown-toc -->

## Data model
### Basic requirements
1. **Data source**: In the first stage the data will primarily be collected from Crunchbase.
2. **Timestamps**: The data will be collected within given time intervals, and one of the motivations of this work is to be able to see a transformation in the dataset throughout a given timeframe.
3. **Mutable API/Data structure**: The data structure can change at any time, new parameters are constantly being added into the Crunchbase API.


## Why NoSQL makes sense for this project?
### Mutable data schema
As described in the previous section,  we rely heavily in an external information source, which means that we do not have much control about how the data is formatted or even the data structure in which the data is collected. Moreover, the structure of the data we collect from Crunchbase can change without notice, since we are not paying for official access to their API. Therefore, the implemented database management system (DBMS) should be able to handle migrations to due changes in the data schema in an easy way.

NoSQL databases are explicitly intended for these scenarios in which there is no guarantee of a fix data schema, and the data is highly mutable.

> While NoSQL databases are equipped to handle normalized data and they are able to sort data within a predefined schema, their respective data models usually allow for **far greater flexibility** than the rigid structure imposed by relational databases. Because of this, NoSQL databases have a reputation for being a **better choice for storing semi-structured and unstructured data**. With that in mind, though, because NoSQL databases don’t come with a predefined schema that often means it’s up to the database administrator to define how the data should be organized and accessed in whatever way makes the most sense for their application.

Taken from [Digital Ocean's: A Comparison of NoSQL Database Management Systems and Models](https://www.digitalocean.com/community/tutorials/a-comparison-of-nosql-database-management-systems-and-models).

A traditional relational database management system (like MySQL, MariaDB or PostgreSQL) would have to be constantly mantained and migrated, due to unexpected changes in the data schema coming from Crunchbase, e.g. any time that Crunchbase adds a new column with a new parameter in their UI, this would most likely create problems in a relational database management system (RDMS) , since the original data schema has been changed. A data engineer would have to manually change the data schema of the RDMS, and perform a _migration_ of the existent information in the RDMS into a database with the new schema. Needless to say, this is a very cumbersome process. 

NoSQL databases, specially **document-oriented databases**, like MongoDB, do not work with a fix data schema and do not need to be mantained every time the data format in Crunchbase is modified. Moreover, they natively work with **JSON files**, so they can directly add the data gathered from Crunchbase without much modifications needed.

## Document-oriented databases
Quite a good description of the benefits of using documented-oriented databases, instead of relational databases. Specially when using data without a rigid data schema.
>In a relational database, [it is not possible or very cumbersome to store objects with a different data schema in the same table]. This is not the case with document databases, which offer the freedom to save multiple documents with different schemas together with no changes to the database itself.
In document databases, documents are not only self-describing but also their schema is dynamic, which means that you **don’t have to define it before you start saving data**. Fields can differ between different documents in the same database, and you can modify the document’s structure at will, adding or removing fields as you go. Documents can be also nested — meaning that a field within one document can have a value consisting of another document — making it possible to store complex data within a single document entry.

>All these attributes make it intuitive to work with document databases from the developer’s perspective. The database facilitates storing actual objects describing data within the application, encouraging experimentation and allowing great flexibility when reshaping data as the software grows and evolves.

Taken from [Digital Ocean's Introduction to document-oriented databases](https://www.digitalocean.com/community/conceptual_articles/an-introduction-to-document-oriented-databases).

### Main characteristics of document-oriented databases
>* **Flexibility and adaptability**: with a high level of control over the data structure, document databases enable experimentation and adaptation to new emerging requirements. New fields can be added right away and existing ones can be changed any time. **It’s up to the developer to decide whether old documents must be amended or the change can be implemented only going forward**.

A very important point for our implementation is mention in the first characteristic. If the data schema of Crunchbase changes we can then very quickly adapt the data pipeline, if we are using a document-oriented database.

>* **Ability to manage structured and unstructured data**: relational databases are well suited for storing data that conforms to a rigid structure. Document databases can be used to handle structured data as well, but they’re also quite useful for storing unstructured data where necessary. 

Taken from [Digital Ocean's Introduction to document-oriented databases](https://www.digitalocean.com/community/conceptual_articles/an-introduction-to-document-oriented-databases).

## Time-series databases
Time-series databases (TSDB) are explicitly design to focus on managing time-dependent datasets, which would be a good fit to fulfil our #2 requirement.

>Time series databases have key architectural design properties that make them very different from other databases. These include **time-stamp data storage and compression**, data lifecycle management, [...].
With a time series database, it is common to request a summary of data over a large time period. This requires going over a range of data points to perform some computation like a percentile increase this month of a metric over the same period in the last six months, summarized by month. This kind of workload is very difficult to optimize for with a distributed key value store. TSDB’s are optimized for exactly this use case giving millisecond level query times over months of data. 

Taken from [InfluxDB: Time-series databases explained](https://www.influxdata.com/time-series-database/) (InfluxDB is currently the most popular time-series database).

### Arguments against TSDBs
A TSDB might be a bit overkill, since most of them tend to be designed towards high-volume/high-frequency data collection. Which is not the case in our project. For example, InfluxDB offers a precision of nanoseconds, if it is required by the task. Although a TSDB would handle the compression of our timestamps better than a document-oriented database, they are normally not designed to handle very mutable data schemas, so they would need more maintenance that a document-oriented DB every time that Crunchbase changes its data schema.

## MongoDB
### Collections
>Generally, having a **large number of collections has no significant performance penalty and results in very good performance**. Distinct collections are very important for high-throughput batch processing.[[Reference](https://www.mongodb.com/docs/manual/core/data-model-operations/)]

### Time-series collections
MongoDB supports creating [_time-series collections_](https://www.mongodb.com/docs/manual/core/timeseries-collections/). Nonetheless, it seems like it was designed for high-frequency operations, since the _granularity_ field that describes the time-series collection accepts the values `seconds`, `minutes` or `hours`, but not days. In our implementation, the data collection is intended to take place with not such a high frequency. So using time-series collections might provide an unncessary pre-optimization of the data handling.

>When you query time series collections, you operate on one document per measurement. Queries on time series collections take advantage of the optimized internal storage format and return results faster.[[Reference](https://www.mongodb.com/docs/manual/core/timeseries-collections/)]

>The implementation of time series collections uses internal collections that reduce disk usage and improve query efficiency. Time series collections automatically order and index data by time.[[Reference](https://www.mongodb.com/docs/manual/core/timeseries-collections/)]


### Indexing
Indexes can be used to improve performance for common queries. Build indexes on fields that appear often in queries and for all operations that return sorted results. Each index requires at least 8 kB of data space.

* **Adding an index has some negative performance impact for write operations**. For collections with high write-to-read ratio, indexes are expensive since each insert must also update any indexes.

* Collections with high read-to-write ratio often benefit from additional indexes. Indexes do not affect un-indexed read operations. When active, each index consumes disk space and memory. 

References: 
* [Operational factors and data model](https://www.mongodb.com/docs/manual/core/data-model-operations/)
* [MongoDB documentation on indexes](https://www.mongodb.com/docs/manual/indexes/)

>The best indexes for your application must take a number of factors into account, **including the kinds of queries you expect, the ratio of reads to writes, and the amount of free memory on your system**.
>When developing your indexing strategy you should have a deep understanding of your application's queries. Before you build indexes, map out the types of queries you will run so that you can build indexes that reference those fields. Indexes come with a performance cost, but are more than worth the cost for frequent queries on large data sets. Consider the relative frequency of each query in the application and whether the query justifies an index.

Taken from: [Indexing strategies](https://www.mongodb.com/docs/manual/applications/indexes/)

### Schema versioning
In order to properly handle changes in the data schema of the pages being scraped a [data model for schema versioning](https://www.mongodb.com/docs/manual/tutorial/model-data-for-schema-versioning/) can be used. So that the application can easily ackowledge a schema change in the data without downtime.

> Using the `schema_version` field, application code can support any number of schema iterations in the same collection by adding dedicated handler functions to the code [without needing downtime or migrations after each schema change].
