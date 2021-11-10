# Telegram Chatbot

This chatbot tracks your spendings in a given month and lists the spent amount of you ask for it.
Currently new spendings can be sent via text messages, however uploading receipts to automatically ingest spendings is planned.

The bot supports a single chat only (as the DB does not have different users) + the spendings of other users are free text in the DB.

## Supported commands:
* ```/list [MM.YYYY]``` - shows all tracked spendings for the given month. If left blank the current month is returned.
* ```/file [MM.YYYY]``` - **NOT IMPLEMENTED YET** returns a file containing all tracked spendings for the given month. If left blank the current month is returned.
* ```/edit <ID> (NAME | AMOUNT)``` - updates an existing spending. You can either change the name or the amount spent
* ```/delete <ID>```- deletes the record
* ```/configure key value```- sets or updates a config value. Currently the only config key in use is: *defaultCurrency*
* ```/show``` - lists all active config values
