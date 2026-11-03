---
layout: post
title: DB/SQL history since college
tags [tech]
date: 2026-11-3
---
This is the endpoint of a late nitght discussion with chatgpt about the evolution of sql since I was in the industry with fun links at the end. Not my own writing. 

# Twenty Years of Database Wars: How SQL Survived NoSQL

*From Oracle and MySQL to Bigtable, MongoDB, and PostgreSQL's unexpected dominance.*

In 2005, the database world seemed relatively settled. Relational databases dominated, SQL was the universal language, and choosing a database largely meant choosing between expensive commercial products and increasingly capable open-source alternatives.

Twenty years later, we've witnessed an explosion of new database architectures, the rise of NoSQL, globally distributed transactions, and the emergence of cloud-native storage. Yet the biggest winner might be a database that was already nearly a decade old in 2005: PostgreSQL.

How did we get here?

## 2000–2005: Open Source Challenges the Database Giants

For decades, enterprise databases were dominated by Oracle, IBM DB2, and Microsoft SQL Server. These were sophisticated systems, but commercial licensing could be extraordinarily expensive.

The internet changed the economics. MySQL, first released in 1995, became a foundation of the LAMP stack (Linux, Apache, MySQL, PHP). It was free, relatively simple to administer, and fast enough for the websites being built by an emerging generation of developers.

PostgreSQL, released under that name in 1996, offered a more ambitious alternative, emphasizing SQL standards, advanced data types, transactions, and extensibility. SQLite, released in 2000, demonstrated that an entire transactional SQL engine could be embedded in an application without a database server.

By 2005, open-source relational databases had become credible alternatives to commercial products. SQL's dominance seemed secure.

But trouble was brewing at Google and Amazon.

## 2005–2010: Google and Amazon Reinvent Storage

As Google scaled its services across thousands of servers, traditional relational databases presented difficulties. Scaling writes generally meant buying bigger machines or manually partitioning data across independent database servers.

Google's 2006 [Bigtable paper](https://research.google/pubs/bigtable-a-distributed-storage-system-for-structured-data/) described a different architecture: an enormous, sparse, sorted map distributed across machines, automatically dividing data into manageable ranges called tablets.

Bigtable abandoned relational queries and general-purpose transactions in favor of scalability and predictable performance.

Amazon's 2007 [Dynamo paper](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf) approached the problem differently. Its distributed key-value system prioritized availability, allowing operations to continue through server failures and network partitions, with conflicts resolved when necessary.

These papers inspired an entire generation of databases. Apache HBase followed Bigtable's model. Cassandra combined Bigtable-like data structures with Dynamo-inspired distribution. Riak explored highly available key-value storage.

The underlying challenge wasn't that SQL was fundamentally wrong. It was that the prevailing implementations of relational databases weren't designed to distribute writes automatically across hundreds of machines.

Nevertheless, the emerging movement became known as **NoSQL**.

## 2009–2015: The NoSQL Gold Rush

MongoDB, introduced in 2009, brought NoSQL to ordinary application developers.

Instead of managing relational tables, developers stored JSON-like documents. Schemas could evolve without traditional migrations, and application objects mapped naturally to database records.

MongoDB's appeal was as much about developer productivity as scalability. During this period, JSON APIs, JavaScript frameworks, and rapidly changing web applications made document databases particularly attractive.

Other systems flourished. Redis made in-memory key-value storage ubiquitous. CouchDB offered document replication. Cassandra promised distributed write scalability. Amazon launched DynamoDB in 2012 as a managed NoSQL service, while Microsoft's Azure DocumentDB, later Cosmos DB, expanded the cloud database market.

For a while, it seemed relational databases might become legacy technology.

But NoSQL came with tradeoffs.

Applications frequently needed relationships between records, transactions spanning multiple objects, and queries that hadn't been anticipated when the data model was designed. Without joins or sufficiently expressive transactions, developers often reimplemented database functionality in application code.

And the promised simplicity of schema-free development could become difficult to maintain as applications matured.

Meanwhile, relational databases weren't standing still.

## 2012–2020: SQL Strikes Back

Two developments helped change the direction of database technology.

The first was Google's Spanner, described in a landmark 2012 paper. Spanner demonstrated that globally distributed databases could provide strong transactional guarantees rather than requiring applications to accept eventual consistency.

This helped inspire a new generation of distributed SQL systems, including CockroachDB, TiDB, and YugabyteDB.

These systems borrowed heavily from NoSQL architectures—automatic partitioning, distributed replication, and consensus protocols—but retained SQL and transactions.

The second, perhaps more consequential development for everyday developers, happened inside PostgreSQL.

In 2014, PostgreSQL 9.4 introduced JSONB, allowing efficient, indexable storage of JSON documents alongside traditional relational data.

Suddenly, developers didn't necessarily have to choose between a document database and a relational database. PostgreSQL could provide both.

PostgreSQL also benefited from an unusually extensible architecture. PostGIS made it a powerful geospatial database. TimescaleDB expanded its usefulness for time-series workloads. Citus introduced distributed sharding capabilities.

Rather than being displaced by specialized databases, PostgreSQL increasingly absorbed their functionality.

## 2020–2026: PostgreSQL Becomes the Default

By the 2020s, PostgreSQL had become the preferred general-purpose database for an enormous range of new applications.

In the 2026 Stack Overflow Developer Survey, 57.9% of respondents reported working with PostgreSQL, versus 35.1% for MySQL and 19.0% for MongoDB. These are developer usage figures, not measures of database installations or revenue, but they illustrate the change in sentiment.

Several factors explain PostgreSQL's success.

First, it is genuinely capable. Its query optimizer, transactional guarantees, indexing options, and extensibility make it a strong default for most applications.

Second, its ecosystem has become enormous. Nearly every cloud provider offers managed PostgreSQL, virtually every development framework supports it, and companies can change hosting providers without necessarily changing databases.

Third, PostgreSQL continues to assimilate competing technologies. The emergence of pgvector brought vector similarity search into ordinary SQL applications, allowing developers to experiment with AI workloads without adopting another database.

But PostgreSQL's popularity conceals some architectural weaknesses.

It still generally depends on a single writable primary. Read replicas can distribute queries, but distributing writes requires sharding, extensions, or independently managed database instances. Its MVCC implementation also creates dead tuples that must be cleaned up through VACUUM.

For individual applications, modern hardware is powerful enough that these limitations may never matter. For platforms hosting hundreds of applications, shared write bottlenecks and resource contention are less attractive.

The rise of PostgreSQL-compatible distributed databases suggests that developers increasingly value **PostgreSQL's interface and ecosystem independently of its underlying storage architecture**.

## The Next Chapter: Distributed SQL, Analytics, and Separation of Storage

While PostgreSQL was winning developer mindshare, several other important developments were unfolding.

**Distributed SQL** systems such as CockroachDB, YugabyteDB, TiDB, and Google's Spanner attempt to preserve SQL and strong transactions while distributing data and write ownership across machines. Their architectures are attractive, although distributed coordination brings latency, operational complexity, and cost.

**Analytical databases** have followed a different path. Snowflake, BigQuery, ClickHouse, and DuckDB demonstrate the advantages of column-oriented storage and specialized execution engines. DuckDB is particularly interesting because it brings high-performance analytical SQL to a single embedded process, much as SQLite did for transactions.

**Storage and compute are separating.** Cloud-native architectures increasingly place durable data in shared or object-based storage while allowing compute resources to scale independently. Systems such as Neon and Amazon Aurora illustrate different approaches to decoupling database compute from storage.

**Open-source licensing is becoming a battleground.** Companies including MongoDB, Elastic, Redis, and Cockroach Labs have moved important products away from conventional permissive open-source licensing. PostgreSQL's permissive license and independent development community have consequently become competitive advantages of their own.

And there's a surprising continuity underneath all these changes: many modern databases, including distributed and analytical systems, still expose SQL.

## Did SQL Win?

The last twenty years of database development weren't simply a story of SQL defeating NoSQL.

NoSQL solved real problems. Bigtable and Dynamo fundamentally changed how distributed storage systems are designed. Cassandra, MongoDB, DynamoDB, and Cosmos DB continue to serve workloads for which their architectures make sense.

But NoSQL also revealed just how valuable relational queries and transactions were. Developers discovered that eliminating database complexity frequently meant moving that complexity into applications.

PostgreSQL's resurgence illustrates another lesson: a mature, extensible platform can outlast specialized competitors by incorporating their best features.

The next generation of databases may distribute writes automatically, separate compute from storage, or use entirely different storage engines. But many will still speak SQL—and an increasing number will speak PostgreSQL's dialect of it.

**Perhaps the real winner of the database wars wasn't PostgreSQL itself, but the relational abstraction it helped preserve.**

---

### Further reading

1. Michael Stonebraker and Andy Pavlo, [*What Goes Around Comes Around... And Around...*](https://db.cs.cmu.edu/papers/2024/whatgoesaround-sigmodrec2024.pdf), 2024.
2. Fay Chang et al., [*Bigtable: A Distributed Storage System for Structured Data*](https://research.google/pubs/bigtable-a-distributed-storage-system-for-structured-data/), 2006.
3. Giuseppe DeCandia et al., [*Dynamo: Amazon's Highly Available Key-value Store*](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf), 2007.
4. James C. Corbett et al., [*Spanner: Google's Globally-Distributed Database*](https://research.google/pubs/spanner-googles-globally-distributed-database/), 2012.
5. Andy Pavlo and Bohan Zhang, [*The Part of PostgreSQL We Hate the Most*](https://www.cs.cmu.edu/~pavlo/blog/2023/04/the-part-of-postgresql-we-hate-the-most.html), 2023.
6. [Stack Overflow Developer Survey 2026 — Databases](https://survey.stackoverflow.co/2026/technology/data/databases).
