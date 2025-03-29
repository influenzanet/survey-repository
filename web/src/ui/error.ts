
export const showError = (error: any) => {
    console.log(error);
    const $e = $('#error');
    $e.empty();
    if(typeof(error) == "object") {
        if(error instanceof Error) {
            $e.append('Error <b class="me-1">' + error.name + '</b> ');
            $e.append('<i>' + error.message +'</i>');
            if(error.stack) {
                $e.append(' : ');
                $e.append(error.stack);
            }
        } else {
            $e.append(JSON.stringify(error));
        }
    } else {
        $e.append(error);
    }
    $e.show();
}