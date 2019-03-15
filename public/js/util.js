var urlParams = new URLSearchParams(window.location.search);

if (urlParams.has('err'))
{

    $.notify({
        icon: "notifications",
        message: urlParams.get('err')

    },{
        type: "danger",
        timer: 4000,
        placement: {
            from: "top",
            align: "center"
        }
    });

}