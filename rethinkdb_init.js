r.db('mercury').table('channels').indexCreate('hash')
r.db('mercury').table('messages').indexCreate('created_at')