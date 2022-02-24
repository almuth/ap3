<?php
/* @var $this MembershipController */
/* @var $model MembershipRegistrationForm */
/* @var $form CActiveForm */

Yii::app()->clientScript->registerCssFile(Yii::app()->theme->baseUrl . '/css/foundation-datepicker.css');
Yii::app()->clientScript->registerScriptFile(Yii::app()->theme->baseUrl . '/js/foundation-datepicker.js', CClientScript::POS_HEAD);
Yii::app()->clientScript->registerScriptFile(Yii::app()->theme->baseUrl . '/js/locales/foundation-datepicker.id.js', CClientScript::POS_HEAD);
?>
<div class="row">
    <div class="small-12 column">
        <div data-alert class="alert-box radius" style="display:none">
            <span></span>
            <a href="#" class="close button">&times;</a>
        </div>
    </div>
</div>
<div class="form">
    <?php
    $form = $this->beginWidget('CActiveForm', [
        'id'                   => 'ubah-member-form',
        // Please note: When you enable ajax validation, make sure the corresponding
        // controller action is handling ajax validation correctly.
        // There is a call to performAjaxValidation() commented in generated controller code.
        // See class documentation of CActiveForm for details on this.
        'enableAjaxValidation' => false,
    ]);
    ?>
    <?php // echo $form->errorSummary($model, 'Error: Perbaiki input', null, ['class' => 'panel callout']); 
    ?>

    <div class="row">
        <div class="small-12 large-6 columns">
            <?php echo CHtml::label('No Telp/HP', 'noTelp'); ?>
            <?php echo CHtml::textField('noTelp', $model->nomor_telp, ['size' => 45, 'maxlength' => 45, 'autofocus' => 'autofocus']); ?>
        </div>
        <div class="small-12 large-6 columns">
            <?php echo CHtml::label('Nama Lengkap', 'namaLengkap'); ?>
            <?php echo CHtml::textField('namaLengkap', $model->nama_lengkap, ['size' => 45, 'maxlength' => 45, 'autofocus' => 'autofocus']); ?>
        </div>
    </div>
    <div class="row">
        <div class="small-12 large-6 columns">
            <?php echo CHtml::label('Tanggal Lahir', 'tanggalLahir'); ?>
            <?php echo CHtml::textField('tanggalLahir', $model->tanggal_lahir); ?>
        </div>
        <div class="small-12 large-6 columns">
            <?php echo CHtml::label('Pekerjaan', 'pekerjaan'); ?>
            <?php echo CHtml::textField('pekerjaan', $model->pekerjaan); ?>
        </div>
    </div>
    <div class="row">
        <div class="small-12 columns">
            <?php echo CHtml::label('Alamat', 'alamat'); ?>
            <?php echo CHtml::textField('alamat', $model->alamat); ?>
        </div>
    </div>
    <div class="row">
        <div class="medium-6 columns">
            <?php echo CHtml::label('Keterangan', 'keterangan'); ?>
            <?php echo CHtml::textField('keterangan', $model->keterangan); ?>
        </div>
    </div>

    <div class="row">
        <div class="small-12 columns">
            <?php echo CHtml::submitButton('Update', ['class' => 'tiny bigfont success button right']); ?>
        </div>
    </div>

    <?php $this->endWidget(); ?>
</div>
<script>
    $(function() {
        $('.tanggalan').fdatepicker({
            format: 'dd-mm-yyyy',
            language: 'id'
        });
    });
    $(document).ready(function() {
        $("#ubah-member-form").submit(function(event) {
            $(".alert-box").slideUp();
            var formData = {
                noTelp: $("#MembershipRegistrationForm_noTelp").val(),
                namaLengkap: $("#MembershipRegistrationForm_namaLengkap").val(),
                tanggalLahir: $("#MembershipRegistrationForm_tanggalLahir").val(),
                pekerjaan: $("#MembershipRegistrationForm_pekerjaan").val(),
                alamat: $("#MembershipRegistrationForm_alamat").val(),
                keterangan: $("#MembershipRegistrationForm_keterangan").val(),
            };
            $.ajax({
                type: "POST",
                url: "<?= $this->createUrl('prosesregistrasi') ?>",
                data: formData,
            }).done(function(data) {
                // console.log(data);
                data = JSON.parse(data)
                $(".alert-box").slideDown(500, function() {
                    if (data.statusCode == 200) {
                        $(".alert-box").removeClass("alert");
                        $(".alert-box").addClass("warning");
                        $(".alert-box>span").html("Sukses: " + data.data.msg + ". Nomor: <strong>" + data.data.nomor + "</strong>")
                    } else {
                        $(".alert-box").removeClass("warning");
                        $(".alert-box").addClass("alert");
                        $(".alert-box>span").html(data.statusCode + ":" + data.error.type + ". " + data.error.description)
                    }
                })

            });

            event.preventDefault();
        });
    });
</script>