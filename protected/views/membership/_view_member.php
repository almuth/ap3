<h4><small>Profil</small> Member</h4>
<hr />
<?php
$this->widget('BDetailView', [
    'data'       => $model,
    'attributes' => [
        'nomor',
        'nomor_telp',
        'nama_lengkap',
        'tanggal_lahir',
        'pekerjaan',
        'alamat',
        'keterangan'
    ],
]);
