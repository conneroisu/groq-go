library IEEE;
use IEEE.STD_LOGIC_1164.all;
use IEEE.NUMERIC_STD.all;

entity tb_ones_comp is end tb_ones_comp;

architecture Behavioral of tb_ones_comp is

    constant N : integer := 8;
    constant M : integer := 32;

    signal Input  : std_logic_vector(N-1 downto 0);
    signal Output : std_logic_vector(N-1 downto 0);
    
    component OnesComplementor
        generic (N : integer := 8);
        port (
            Input  : in  std_logic_vector(N-1 downto 0);
            Output : out std_logic_vector(N-1 downto 0)
            );
    end component;

    component OnesComplementor_2
        generic (M : integer := 32);
        port (
            Input2  : in  std_logic_vector(M-1 downto 0);
            Output2 : out std_logic_vector(M-1 downto 0)
            );
    end component;
begin
    DUT0 : OnesComplementor
        generic map (N => N)
        port map (
            Input  => Input,
            Output => Output
            );

    DUT1 : OnesComplementor_2
        generic map (M => M)
        port map (
            Input2 => Input,
            Output2 => Output
            );

    stim_proc : process
    begin
        for i in 0 to 2**N-1 loop
            Input <= std_logic_vector(to_unsigned(i, N));
            wait for 10 ns;
            assert Output = not Input report "End of testbench simulation" severity failure;
        end loop;
        for i in 0 to 2**M-1 loop
            Input <= std_logic_vector(to_unsigned(i, M));
            wait for 10 ns;
            assert Output = not Input report "End of testbench simulation" severity failure;
        end loop;

    end process;

end Behavioral;
