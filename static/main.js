$(function() {
  
  $('#dropzone').on('dragover', function() {
    $(this).addClass('hover');
  });
  
  $('#dropzone').on('dragleave', function() {
    $(this).removeClass('hover');
  });
  
  $('#dropzone input').on('change', function(e) {
    var file = this.files[0];

    $('#dropzone').removeClass('hover');    
    $('#dropzone').addClass('dropped');

    document.getElementById("fileform").submit.click();
    var ext = file.name.split('.').pop();
      
    $('#dropzone div').html(ext);
  });
});