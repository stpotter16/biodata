document.querySelectorAll('tbody tr').forEach(row => {
    row.addEventListener('click', () => {
        window.location.href = `/entry/${row.cells[0].textContent}/edit`;
    });
});
