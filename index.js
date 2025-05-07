const { Client } = require('pg');

async function connectToPrimary(connectionConfig) {
  const client = new Client(connectionConfig);
  try {
    await client.connect();
    const res = await client.query('SELECT * from my_table;');
    console.log('Connected to PostgreSQL!', res.rows[0]);
    const insertRes = await client.query("INSERT INTO my_table (data) VALUES ('Testing Insert - 5')")
    console.log('Insert Status', insertRes.rowCount);
    const tableRes = await client.query('SELECT * from my_table;');
    console.log('Retrieving db: ', tableRes.rows);
  } catch (err) {
    console.error('Error connecting to Primary PostgreSQL', err);
  } finally {
    await client.end();
  }
}

async function connectToReplica(connectionConfig) {
  const client = new Client(connectionConfig);
  try {
    await client.connect();
    const res = await client.query('SELECT * from my_table;');
    console.log('Connected to Replica PostgreSQL!', res.rows);
  } catch (err) {
    console.error('Error connecting to Replica PostgreSQL', err);
  } finally {
    await client.end();
  }
}

// Configuration for the primary server
const primaryConfig = {
  host: 'localhost',
  port: 5432,
  user: 'postgres',
  password: '', // Replace with your password if set
  database: 'test_replication' // Replace with the database you want to connect to
};

// Configuration for the standby server
const standbyConfig = {
  host: 'localhost',
  port: 5433, // If you used a different port for the standby
  user: 'postgres',
  password: '',
  database: 'test_replication'
};

connectToPrimary(primaryConfig);
connectToReplica(standbyConfig); // You can connect to the standby for read operations