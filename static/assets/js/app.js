
  $(function() {
    $('#admin-offcanvas').load('/sidebar.html')
    //$('#_header_').load('/header.html')
    $('#_footer_').load('/footer.html')

    var $fullText = $('.admin-fullText');
    $('#admin-fullscreen').on('click', function() {
      $.AMUI.fullscreen.toggle();
    });

    $(document).on($.AMUI.fullscreen.raw.fullscreenchange, function() {
      $fullText.text($.AMUI.fullscreen.isFullscreen ? '退出全屏' : '开启全屏');
    });


  });

