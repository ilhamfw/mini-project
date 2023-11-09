
Entities: 
"users" entity, 
   - Email/username
   - Password    
   - Deposit amount DEFAULT 0
main entity rental games with mandatory attributes for this entity are:
   - Name
   - availability / stock_availability
   - Rental_costs
   - Category

Table Device 
rental history table.

Saya ingin membuat rental Playstation (PS) yang terdiri dari PS2, PS3, PS4. Setiap user bisa merental PS yang mereka mau dengan cost harian, setiap user harus memiliki akun dan memiliki saldo yang disimpan di field Deposit. Saya ingin setelah user login. User memiliki 3 fitur, User mednapatkan notifikasi email, yaitu top up saldo dan merental Playstation

bisa untuk 

POST (/users/register)
POST (/users/login)

GET (/rent) : menampilkan seluruh yang sedang rental yang diambil dari table rent
GET (/rent/:id) : menampilkan berdasarkan id

POST (/rent) : Membuat Rental yang disimpan di table rent lalu ditambahkan juga ke table history. Disini harus implementasi agar saldo yang dimiliki pengguna bisa untuk membayar sesuai harga rental yang ada. Jika tidak maka user harus top up saldo terlebih dahulu

GET (/history) : Menampilkan seluruh history users yang rental ##

PUT (/deposit) : Menampilkan seluruh history users yang rental

Tambahkan Table top up history
